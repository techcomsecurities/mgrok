# mgrok
`mgrok` package leverages [ngrok](http://ngrok.com) to help your program easily making local tunnel.

## Install
```
$ go get github.com/techcomsecurities/mgrok
```

## Example
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/techcomsecurities/mgrok"
)

func main() {
	var ngrok = "/path/to/your/ngrok"
	var port = "1203"
	http.HandleFunc("/", homeHandle)

	tunnels, err := mgrok.Run(ngrok, "http", port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Your tunnels: ", tunnels)
	log.Println("Inspect request at http://localhost:4040")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Let's Go!")
}
```

### Document
Refer to [godoc](https://godoc.org/github.com/techcomsecurities/mgrok)