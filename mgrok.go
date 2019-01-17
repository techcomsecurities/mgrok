package mgrok

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"os"
	"os/exec"
	"os/signal"
)

var ngrokCmd *exec.Cmd
var control chan os.Signal
var timeout = 10 // in second
var notStartedErr = errors.New("ngrok has not started yet")

func init() {
	control = make(chan os.Signal, 2)
	signal.Notify(control, os.Interrupt, os.Kill)
}

// Run runs ngrok at given path and arguments, 
// returns ngrok tunnels assigned to you, and error if has
func Run(ngrok string, args ...string) (tunnels []string, err error) {
	ngrokCmd = exec.Command(ngrok, args...)
	go ngrokCmd.Run()
	go handleExit()

	for i := 0; i < timeout; i++ {
		time.Sleep(1*time.Second)
		tunnels, err = getTunnels()
		if err == notStartedErr {
			continue
		}
		break
	}

	return tunnels, err
}

// SetTimeout sets time (in second) to wait before sending get tunnels request to ngrok process
// default is 10 seconds
func SetTimeout(t int) {
	timeout = t
}

// Close releases all resouces mgrok uses
// Remember to call it before exit your program
func Close() {
	if ngrokCmd != nil {
		ngrokCmd.Process.Kill()
	}
}

func handleExit() {
	<-control
	Close()
	os.Exit(0)
}

func getTunnels() (tunnels []string, err error) {
	resp, err := http.Get("http://localhost:4040/api/tunnels")
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return tunnels, notStartedErr
		}
		return tunnels, err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tunnels, err
	}
	
	if resp.StatusCode != 200 {
		return tunnels, fmt.Errorf("HTTP status code %d - Response %s", resp.StatusCode, string(body))
	}

	var s = struct{
		Tunnels []struct{
			PublicURL string `json:"public_url"`
		} `json:"tunnels"`
	}{}
	if err := json.Unmarshal(body, &s); err != nil {
		return tunnels, err
	}

	if len(s.Tunnels) == 0 {
		return tunnels, notStartedErr
	}

	for _, t := range s.Tunnels {
		tunnels = append(tunnels, t.PublicURL)
	}
	return tunnels, nil
}