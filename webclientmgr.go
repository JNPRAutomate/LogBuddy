package logbuddy

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/websocket"
)

const (
	//CookieName the cookie key that represents the session
	CookieName = "lbid"
)

//NewWebClientMgr Returns an initalized web client manager
func NewWebClientMgr() *WebClientMgr {
	//generate new key if it does not exist
	return &WebClientMgr{serverMgr: NewServerManager(), ClientServers: make(map[string][]int), Clients: make(map[string]*http.Cookie), sCookie: securecookie.New(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))}
}

//WebClientMgr manages web clients to ensure persistance across sessions
type WebClientMgr struct {
	Clients       map[string]*http.Cookie //maps clients based upon cookies
	ClientServers map[string][]int
	sCookie       *securecookie.SecureCookie
	serverMgr     *ServerManager //serverMgr Interaction with the server manager to review jobs
}

//StartSession generates a new cookie for clients and adds the session to the WebClientMgr
func (wcm *WebClientMgr) StartSession(w http.ResponseWriter, r *http.Request) {
	if wcm.checkSession(r) {
		//existing cookie set
		return
	}
	cookie := wcm.generateCookie()
	wcm.Clients[cookie.Value] = cookie
	http.SetCookie(w, cookie)
}

//StartWSSession Starts a websocket session with the WebClientMgr allows it to setup existing session information
func (wcm *WebClientMgr) StartWSSession(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) (string, []chan LogMessage, []chan CtrlChanMsg) {
	if wcm.checkSession(r) {
		//session exists
		//reconect logging connections
		if cookie, err := r.Cookie(CookieName); err == nil {
			var logChans []chan LogMessage
			var ctrlChans []chan CtrlChanMsg
			if len(wcm.ClientServers[cookie.Value]) > 0 {
				for item := range wcm.ClientServers[cookie.Value] {
					wscm := &WSClientMessage{Type: RestartMsg}
					config := wcm.serverMgr.ServerConfigs[wcm.ClientServers[cookie.Value][item]]
					jsonConfig, _ := config.MarshalJSON()
					wscm.Data = jsonConfig
					conn.SetWriteDeadline(time.Now().Add(writeWait))
					clientMsg, _ := wscm.MarshalJSON()
					if err := conn.WriteMessage(websocket.TextMessage, clientMsg); err != nil {
						return cookie.Value, nil, nil
					}
					logChan, ctrlChan := wcm.ReconnectSession(wcm.ClientServers[cookie.Value][item])
					logChans = append(logChans, logChan)
					ctrlChans = append(ctrlChans, ctrlChan)
				}
			}
			return cookie.Value, logChans, ctrlChans
		}
	}
	cookie := wcm.generateCookie()
	wcm.Clients[cookie.Value] = cookie
	http.SetCookie(w, cookie)
	return cookie.Value, nil, nil
}

//checkSession checks for session exists
func (wcm *WebClientMgr) checkSession(r *http.Request) bool {
	if cookie, err := r.Cookie(CookieName); err == nil {
		if _, ok := wcm.Clients[cookie.Value]; ok {
			return true
		}
	}
	return false
}

//StartServer starts a new server for a web client
func (wcm *WebClientMgr) StartServer(client string, config *ServerConfig) (chan LogMessage, chan CtrlChanMsg) {
	//add server ids to client servers
	//return id, error
	chanID, err := wcm.serverMgr.StartServer(config)
	if err != nil {
		log.Println("Error", err)
		return nil, nil
	}
	wcm.bindServer(client, chanID)
	logChan, ctrlChan := wcm.serverMgr.Register(chanID)
	return logChan, ctrlChan
}

//bindServer Binds a client to a server
func (wcm *WebClientMgr) bindServer(client string, id int) {
	wcm.ClientServers[client] = append(wcm.ClientServers[client], id)
}

//Binds a client to a server
func (wcm *WebClientMgr) unbindServer(client string, id int) {
	wcm.ClientServers[client] = append(wcm.ClientServers[client][:id], wcm.ClientServers[client][id+1:]...)
}

//StopSession stops an existing session for a client. Also stops all existing servers.
func (wcm *WebClientMgr) StopSession(id string) {
	for item := range wcm.ClientServers[id] {
		//stop server
		wcm.serverMgr.StopServer(wcm.ClientServers[id][item])
		//unbind server from client
		wcm.unbindServer(id, wcm.ClientServers[id][item])
	}
	delete(wcm.Clients, id)
}

//ReconnectSession Returns an existing Message channel based on chanID
func (wcm *WebClientMgr) ReconnectSession(chanID int) (chan LogMessage, chan CtrlChanMsg) {
	logChan, ctrlChan := wcm.serverMgr.Register(chanID)
	if logChan == nil || ctrlChan == nil {
		return nil, nil
	}
	return logChan, ctrlChan
}

//generateCookie generates a new cookie for a client
func (wcm *WebClientMgr) generateCookie() *http.Cookie {
	var err error
	var encoded string
	//TODO: Make cookies less predictable
	value := map[string]string{
		"foo": "bar",
	}
	if encoded, err = wcm.sCookie.Encode(CookieName, value); err == nil {
		cookie := &http.Cookie{
			Name:    CookieName,
			Value:   encoded,
			Path:    "/",
			Expires: time.Now().Add(6 * time.Hour),
		}
		return cookie
	}
	panic(err)
}
