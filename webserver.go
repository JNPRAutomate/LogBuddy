package logbuddy

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
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
	Address   string            //Address the address to listen on
	Router    *mux.Router       //Router for handling requests
	ServerMgr *ServerManager    //ServerMgr Interaction with the server manager to review jobs
	wsConns   []*websocket.Conn //Conns all open connection
}

//Listen set webserver to listen
func (ws *WebServer) Listen() error {
	r := mux.NewRouter()
	r.HandleFunc("/", ws.homeHandler)
	r.HandleFunc("/logs", ws.wsServeLogs)
	http.Handle("/", r)
	if err := http.ListenAndServe(ws.Address, nil); err != nil {
		return err
	}
	return nil
}

//Close Stop the web server and close existing clients
func (ws *WebServer) Close() error {
	//stop all connections
	//stop all open websocket connections
	//stop server from listening
	return nil
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
	go func(conn *websocket.Conn) {
		pingTicker := time.NewTicker(pingRate)
		for {
			select {
			case <-pingTicker.C:
				log.Println("PING")
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			}
		}
	}(conn)
	for {
		msgType, data, err := conn.ReadMessage()
		log.Println("MsgType", msgType)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		switch msgType {
		case websocket.BinaryMessage:
			//currently not implemented
			log.Println("Bin")
			break
		case websocket.TextMessage:
			//handle json requests
			log.Println("Text")
			log.Println(string(data))
		case websocket.CloseMessage:
			conn.Close()
			break
		case websocket.PingMessage:
			break
		case websocket.PongMessage:
			break
		default:
			log.Println("READ")
		}
	}
}

func (ws *WebServer) wsOriginChecker(r *http.Request) bool {
	log.Println("ORIGN")
	return true
}

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
