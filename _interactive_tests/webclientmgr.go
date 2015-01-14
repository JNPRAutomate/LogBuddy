package logbuddy

import "net/http"

//WebClientMgr manages web clients to ensure persistance across sessions
type WebClientMgr struct {
	Clients map[http.Cookie]string //maps clients based upon cookies

}

//StartSession generates a new cookie for clients and adds the session to the WebClientMgr
func (wcm *WebClientMgr) StartSession() {

}

//StopSession stops an existing session for a llient. Also stops all existing servers.
func (wcm *WebClientMgr) StopSession() {

}
