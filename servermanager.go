package logbuddy

import (
	"math/rand"
	"time"
)

//ServerManager Manages listening servers
type ServerManager struct {
	CtrlChans map[int]chan CtrlChanMsg
}

//NewServerManager Creates a new server manager with an initalized CtrlChans map
func NewServerManager() *ServerManager {
	return &ServerManager{CtrlChans: make(map[int]chan CtrlChanMsg)}
}

//StartServer Start a new server with a server config
func (s *ServerManager) StartServer(config *ServerConfig) (id int, err error) {
	//set default id to 0
	id = s.getID()

	if config.Type == "tcp4" || config.Type == "tcp6" || config.Type == "tcp" {
		msgChan := make(chan Message)
		ctrlChan := make(chan CtrlChanMsg)
		listener := &TCPServer{Config: config, msgChan: msgChan, ctrlChan: ctrlChan}
		listener.setListener()
		go listener.Listen()
	} else if config.Type == "udp4" || config.Type == "udp6" || config.Type == "udp" {
		msgChan := make(chan Message)
		ctrlChan := make(chan CtrlChanMsg)
		listener := &UDPServer{Config: config, ctrlChan: ctrlChan, msgChan: msgChan}
		listener.setListener()
		go listener.Listen()
	} else {
		//set error for not found
		return id, err
	}
	return id, err
}

//StopServer Stop a server specified by IP
func (s *ServerManager) StopServer(id int) error {
	//stop instance of server based on ID
	if _, ok := s.CtrlChans[id]; ok {
		s.CtrlChans[id] <- CtrlChanMsg{Type: StopMsg}
	}
	return nil
}

func (s *ServerManager) getID() (id int) {
	rand.Seed(time.Now().Unix() * rand.Int63())
	id = rand.Int()

	if _, ok := s.CtrlChans[id]; !ok {
		return id
	}
	return s.getID()
}
