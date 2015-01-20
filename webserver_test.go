package logbuddy

//go:generate go-bindata -pkg "logbuddy" -o logbuddy_bindata.go static/...

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestBasicWebServer(t *testing.T) {
	testServer := "localhost:8080"
	ws := &WebServer{Address: testServer}
	go ws.Listen()
	defer ws.Close()
	time.Sleep(1 * time.Second)
	res, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	t.Log("Data Bytes Recieved:", len(data))
	res.Body.Close()
}

func TestBasicWebSocketServer(t *testing.T) {
	testServer := "localhost:8081"
	ws := &WebServer{Address: testServer}
	go ws.Listen()
	time.Sleep(1 * time.Second)
	wsClient := &websocket.Dialer{HandshakeTimeout: 3 * time.Second,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
	clientHeader := http.Header{}
	clientHeader.Add("Origin", "http://localhost:8081")
	clientHeader.Add("Host", testServer)
	clientHeader.Add("Sec-WebSocket-Extensions", "perMessage-deflate; client_max_window_bits")
	cConn, _, err := wsClient.Dial("ws://localhost:8081/logs", clientHeader)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	err = cConn.WriteMessage(websocket.TextMessage, []byte("Hello there"))
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	cConn.Close()
}
