package main

import "github.com/jnprautomate/logbuddy"

func main() {
	testServer := "localhost:8080"
	ws := &logbuddy.WebServer{Address: testServer}
	ws.Listen()
	defer ws.Close()
}
