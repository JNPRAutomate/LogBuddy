package main

//go:generate go-bindata -pkg "logbuddy" -o ../logbuddy_bindata.go ../static/...

import "github.com/jnprautomate/logbuddy"

func main() {
	testServer := "localhost:8080"
	ws := &logbuddy.WebServer{Address: testServer}
	ws.Listen()
	defer ws.Close()
}
