package logbuddy

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

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
	//check for existing session
	clientID, logChans := ws.ClientMgr.StartWSSession(w, r, conn)

	for item := range logChans {
		go func(logChan chan LogMessage) {
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
		}(logChans[item])
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
				//TODO: Add return minding message based upon cookie
				cm := &ClientMessage{}
				if err := json.Unmarshal(data, cm); err != nil {
					//error in decoding JSON
					conn.SetWriteDeadline(time.Now().Add(writeWait))
					if err := conn.WriteMessage(websocket.TextMessage, []byte("JSON Error")); err != nil {
						return
					}
				}
				//process message
				logChan := ws.ClientMgr.StartServer(clientID, &cm.ServerConfig)
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
	//TODO: Check that origin is of the same page
	return true
}

//wsError Handles errors for WebSocket connections
func (ws *WebServer) wsError(w http.ResponseWriter, r *http.Request, status int, reason error) {
	log.Println(status, reason)
}

//RegisterLogger Registers a logger to be sent to the connection
func (ws *WebServer) RegisterLogger(id int) (msgChan chan LogMessage, err error) {
	return nil, nil
}
