package logbuddy

import (
	"net"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	// Time allowed to write the file to the client.
	writeWait = 15 * time.Second
	// Max time to wait for the next pong message
	pongKeepAlive = 60 * time.Second
	// Rate to send ping messages to client
	pingRate = (pongKeepAlive * 9) / 10
)

//WebServer Serves the user front end and APIs
type WebServer struct {
	listener  net.Listener      //TCP listener
	Address   string            //Address the address to listen on
	wsConns   []*websocket.Conn //Conns all open connection
	ClientMgr *WebClientMgr     //Client manager manages the state of clients
}

//Listen set webserver to listen
func (ws *WebServer) Listen() error {
	var err error
	r := mux.NewRouter()
	ws.ClientMgr = NewWebClientMgr()
	r.HandleFunc("/", ws.HomeHandler).Methods("GET")
	r.HandleFunc("/logs", ws.wsServeLogs)
	r.HandleFunc("/static/{file:[a-zA-Z/.-]+}", ws.ServeStatic).Methods("GET")
	addr, err := net.ResolveTCPAddr("tcp", ws.Address)
	if err != nil {
		return err
	}
	ws.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	if err = http.Serve(ws.listener, r); err != nil {
		return err
	}
	return nil
}

//Close Stop the web server and close existing clients
func (ws *WebServer) Close() error {
	//stop all connections
	//stop all open websocket connections
	//stop server from listening
	for item := range ws.wsConns {
		if ws.wsConns[item] != nil {
			ws.wsConns[item].SetReadDeadline(time.Now())
			ws.wsConns[item].Close()
		}

	}
	return ws.listener.Close()
}

//ServeStatic Serves static content
func (ws *WebServer) ServeStatic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["file"]
	data, err := Asset(strings.Join([]string{"static/", fileName}, ""))
	if err != nil {
		// Asset was not found.
		http.NotFound(w, r)
	}
	_, file := path.Split(fileName)
	fType := strings.Split(file, ".")
	if len(fType) == 2 {
		if fType[1] == "js" {
			w.Header().Set("Content-Type", "text/javascript")
		} else if fType[1] == "css" {
			w.Header().Set("Content-Type", "text/css")
		}
	}
	if len(fType) == 3 {
		if fType[2] == "js" {
			w.Header().Set("Content-Type", "text/javascript")
		} else if fType[2] == "css" {
			w.Header().Set("Content-Type", "text/css")
		}
	}
	w.Write(data)
}

const homeHTML = `<!DOCTYPE html>
<html lang="en">
	<head>
		<title>Basic Logs</title>
	</head>
	<body>
		<div id="logData"></div>
		<script type="text/javascript">
		(function() {
			var data = document.getElementById("logData");
			var conn = new WebSocket("ws://{{.Host}}/logs");
			conn.onopen = function(evt) {
				console.log("MSG sent")
				conn.send("hello THERE");
				data.textContent = 'Connection Open';
			}
			conn.onclose = function(evt) {
				data.textContent = 'Connection closed';
			}
			conn.onmessage = function(evt) {
				console.log(evt);
				data.textContent = evt.data;
			}
			;
			})();
		</script>
	</body>
</html>
`
