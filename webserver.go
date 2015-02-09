package logbuddy

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	// Time allowed to write the file to the client.
	writeWait = 15 * time.Second
	// Max time to wait for the next pong LogMessage
	pongKeepAlive = 60 * time.Second
	// Rate to send ping LogMessages to client
	pingRate = (pongKeepAlive * 9) / 10
)

//WebServer Serves the user front end and APIs
type WebServer struct {
	listener  net.Listener      //TCP listener
	Address   string            //Address the address to listen on
	wsConns   []*websocket.Conn //Conns all open connection
	ClientMgr *WebClientMgr     //Client manager manages the state of clients
	wg        sync.WaitGroup
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
