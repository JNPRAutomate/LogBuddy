package logbuddy

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

const (
	//CookieName the cookie key that represents the session
	CookieName = "lbid"
)

//NewWebClientMgr Returns an initalized web client manager
func NewWebClientMgr() *WebClientMgr {
	//generate new key if it does not exist
	return &WebClientMgr{ServerMgr: NewServerManager(), ClientServers: make(map[string][]int), Clients: make(map[string]*http.Cookie), sCookie: securecookie.New(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))}
}

//WebClientMgr manages web clients to ensure persistance across sessions
type WebClientMgr struct {
	Clients       map[string]*http.Cookie //maps clients based upon cookies
	ClientServers map[string][]int
	sCookie       *securecookie.SecureCookie
	ServerMgr     *ServerManager //ServerMgr Interaction with the server manager to review jobs
}

//StartSession generates a new cookie for clients and adds the session to the WebClientMgr
func (wcm *WebClientMgr) StartSession(w http.ResponseWriter, r *http.Request) {
	if wcm.checkSession(r) {
		//session exists
		//reconect logging connections

	} else {
		//set new cookie
		cookie := wcm.generateCookie()
		wcm.Clients[cookie.Value] = cookie
		http.SetCookie(w, cookie)
	}

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
func (wcm *WebClientMgr) StartServer(id string, config *ServerConfig) chan Message {
	//add server ids to client servers
	//return id, error
	chanID, err := wcm.ServerMgr.StartServer(config)
	if err != nil {
		log.Println("Error", err)
		return nil
	}
	logChan := wcm.ServerMgr.Register(chanID)
	return logChan
}

//BindServer Binds a client to a server
func (wcm *WebClientMgr) BindServer(client string, id int) {
	wcm.ClientServers[client] = append(wcm.ClientServers[client], id)
}

//StopSession stops an existing session for a llient. Also stops all existing servers.
func (wcm *WebClientMgr) StopSession(id string) {
	delete(wcm.Clients, id)
	//TODO: loop and stop all associated servers
}

func (wcm *WebClientMgr) ReconnectSession(chanID int) chan Message {
	logChan := wcm.ServerMgr.Register(chanID)
	if logChan == nil {
		return nil
	}
	return logChan
}

//generateCookie generates a new cookie for a client
func (wcm *WebClientMgr) generateCookie() *http.Cookie {
	var err error
	var encoded string
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
