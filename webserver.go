package logbuddy

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"text/template"
	"time"
)

var (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second
	// Max time to wait for the next pong message
	pongKeepAlive = 60 * time.Second
	// Rate to send ping messages to client
	pingRate = (pongKeepAlive * 9) / 10
)

//WebServer Serves the user front end and APIs
type WebServer struct {
	listener  net.Listener      //TCP listener
	Address   string            //Address the address to listen on
	Router    *mux.Router       //Router for handling requests
	ServerMgr *ServerManager    //ServerMgr Interaction with the server manager to review jobs
	wsConns   []*websocket.Conn //Conns all open connection
}

//Listen set webserver to listen
func (ws *WebServer) Listen() error {
	var err error
	r := mux.NewRouter()
	r.HandleFunc("/", ws.homeHandler)
	r.HandleFunc("/logs", ws.wsServeLogs)
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

func (ws *WebServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	home := template.Must(template.New("").Parse(homeHTML))
	var homeData = struct {
		Host string
	}{
		r.Host,
	}
	home.Execute(w, &homeData)
}

func (ws *WebServer) staticHandler(w http.ResponseWriter, r *http.Request) {

}

func (ws *WebServer) wsServeLogs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 2 * time.Second,
		CheckOrigin:      ws.wsOriginChecker,
		Error:            ws.wsError}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	//handle websocket connection
	ws.wsConns = append(ws.wsConns, conn)
	conn.SetReadLimit(1024)
	conn.SetReadDeadline(time.Now().Add(pongKeepAlive))
	conn.SetPongHandler(func(string) error {
		log.Println("PONG")
		conn.SetReadDeadline(time.Now().Add(pongKeepAlive))
		return nil
	})
	//Send WebSocket PING messages
	go func(conn *websocket.Conn) {
		pingTicker := time.NewTicker(pingRate)
		msgTicker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-pingTicker.C:
				log.Println("PING")
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			case <-msgTicker.C:
				log.Println("TESTMSG")
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.TextMessage, []byte("TESTMSG")); err != nil {
					return
				}
			}
		}
	}(conn)
	go func(conn *websocket.Conn) {
		for {
			msgType, data, err := conn.ReadMessage()
			log.Println("MsgType", msgType)
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
			//Handle various messages
			switch msgType {
			//Handle text messages
			case websocket.TextMessage:
				//handle json requests
				log.Println("Text")
				log.Println(string(data))
				//handle msg resistrations
			//Handle binary messages
			case websocket.BinaryMessage:
				//currently not used
				log.Println("Bin")
			//Handle close messages
			case websocket.CloseMessage:
				//Closing connection
				conn.Close()
				break
			}
		}
	}(conn)
}

//wsOriginChecker Checks the origin request and validates the request
func (ws *WebServer) wsOriginChecker(r *http.Request) bool {
	return true
}

//wsError Handles errors for WebSocket connections
func (ws *WebServer) wsError(w http.ResponseWriter, r *http.Request, status int, reason error) {
	log.Println(status, reason)
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
