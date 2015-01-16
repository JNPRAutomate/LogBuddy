package logbuddy

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

func NewWebClientMgr() *WebClientMgr {
	//generate new key if it does not exist
	return &WebClientMgr{ClientServers: make(map[string][]int), Clients: make(map[string]*http.Cookie), sCookie: securecookie.New(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))}
}

//WebClientMgr manages web clients to ensure persistance across sessions
type WebClientMgr struct {
	Clients       map[string]*http.Cookie //maps clients based upon cookies
	ClientServers map[string][]int
	sCookie       *securecookie.SecureCookie
	ServerMgr     *ServerManager //ServerMgr Interaction with the server manager to review jobs
}
//TODO: Add in server manager here. It makes the most sense as it is bound to the web client

//StartSession generates a new cookie for clients and adds the session to the WebClientMgr
func (wcm *WebClientMgr) StartSession(w http.ResponseWriter) {
	cookie := wcm.generateCookie()
	wcm.Clients[cookie.Value] = cookie
	http.SetCookie(w, cookie)
}

func (wcm *WebClientMgr) BindServer(client string, id int) {
	wcm.ClientServers[client] = append(wcm.ClientServers[client], id)
}

//StopSession stops an existing session for a llient. Also stops all existing servers.
func (wcm *WebClientMgr) StopSession(id string) {
	delete(wcm.Clients, id)
	//TODO: loop and stop all associated servers
}

//generateCookie generates a new cookie for a client
func (wcm *WebClientMgr) generateCookie() *http.Cookie {
	var err error
	var encoded string
	value := map[string]string{
		"foo": "bar",
	}
	if encoded, err = wcm.sCookie.Encode("lbid", value); err == nil {
		cookie := &http.Cookie{
			Name:    "lbid",
			Value:   encoded,
			Path:    "/",
			Expires: time.Now().Add(6 * time.Hour),
		}
		return cookie
	}
	panic(err)
}
