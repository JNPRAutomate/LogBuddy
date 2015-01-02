package logbuddy

import (
	"testing"
)

func TestBasicWebServer(t *testing.T) {
	ws := &WebServer{Address: ":8080"}
	ws.Listen()
}
