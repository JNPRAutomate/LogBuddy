package logbuddy

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
	ServerMgr *ServerManager    //ServerMgr Interaction with the server manager to review jobs
	wsConns   []*websocket.Conn //Conns all open connection
}

//Listen set webserver to listen
func (ws *WebServer) Listen() error {
	var err error
	r := mux.NewRouter()
	ws.ServerMgr = NewServerManager()
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

func (ws *WebServer) wsServeLogs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 2 * time.Second,
		CheckOrigin:      ws.wsOriginChecker,
		Error:            ws.wsError}
	//
	//
	// NEED TO IMPLEMENT COOKIES.
	// THIS WILL BIND THE USER SESSION TO ALL THEIR WEBSOCKETS AND PORTS THAT ARE OPEN.
	//
	//
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
	//Recieve WebSocket Messages
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
				//handle different message requests here
				cm := &ClientMessage{}
				if err := json.Unmarshal(data, cm); err != nil {
					//error in decoding JSON
					conn.SetWriteDeadline(time.Now().Add(writeWait))
					if err := conn.WriteMessage(websocket.TextMessage, []byte("JSON Error")); err != nil {
						return
					}
				}
				//process message
				chanID, err := ws.ServerMgr.StartServer(&cm.ServerConfig)
				if err != nil {
					log.Println("Error", err)
				}
				logChan := ws.ServerMgr.Register(chanID)
				if logChan == nil {
					//channel not found
					conn.SetWriteDeadline(time.Now().Add(writeWait))
					if err := conn.WriteMessage(websocket.TextMessage, []byte("NOT FOUND")); err != nil {
						return
					}
				}
				go func() {
					for {
						select {
						case m := <-logChan:
							conn.SetWriteDeadline(time.Now().Add(writeWait))
							jsonMsg, _ := m.MarshalJSON()
							if err := conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
								return
							}
						}
					}
				}()
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

//RegisterLogger Registers a logger to be sent to the connection
func (ws *WebServer) RegisterLogger(id int) (msgChan chan Message, err error) {
	return nil, nil
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
