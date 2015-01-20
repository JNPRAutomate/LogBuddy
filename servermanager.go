package logbuddy

import (
	"errors"
	"math/rand"
	"time"
)

//ServerManager Manages listening servers
type ServerManager struct {
	ServerConfigs map[int]*ServerConfig
	CtrlChans     map[int]chan CtrlChanMsg
	MsgChans      map[int]chan Message
	ErrChans      map[int]chan error
	msgRouter     *MsgRouter
}

//NewServerManager Creates a new server manager with an initalized CtrlChans map
func NewServerManager() *ServerManager {
	return &ServerManager{CtrlChans: make(map[int]chan CtrlChanMsg),
		MsgChans:      make(map[int]chan Message),
		msgRouter:     &MsgRouter{},
		ErrChans:      make(map[int]chan error),
		ServerConfigs: make(map[int]*ServerConfig)}
}

//StartServer Start a new server with a server config
func (s *ServerManager) StartServer(config *ServerConfig) (id int, err error) {
	//Get a new unused ID
	id = s.getID()
	config.ID = id
	//CHECK FOR LISTENING PORT
	if config.Type == "tcp4" || config.Type == "tcp6" || config.Type == "tcp" {
		s.MsgChans[id] = make(chan Message)
		s.CtrlChans[id] = make(chan CtrlChanMsg)
		s.ErrChans[id] = make(chan error)
		listener := &TCPServer{Config: config, msgChan: s.MsgChans[id], ctrlChan: s.CtrlChans[id], errChan: s.ErrChans[id]}
		s.ServerConfigs[id] = config
		go s.handleCtrl(id)
		go s.handleErrors(id)
		listener.setListener()
		go listener.Listen()
		return id, nil
	} else if config.Type == "udp4" || config.Type == "udp6" || config.Type == "udp" {
		s.MsgChans[id] = make(chan Message)
		s.CtrlChans[id] = make(chan CtrlChanMsg)
		s.ErrChans[id] = make(chan error)
		s.ServerConfigs[id] = config
		listener := &UDPServer{Config: config, msgChan: s.MsgChans[id], ctrlChan: s.CtrlChans[id], errChan: s.ErrChans[id]}
		go s.handleErrors(id)
		listener.setListener()
		go listener.Listen()
		return id, nil
	}
	return id, errors.New("Unable to create specified server")
}

func (s *ServerManager) handleCtrl(id int) {
	for {
		select {
		case msg := <-s.CtrlChans[id]:
			if msg.Type != BadMsg {
				//	s.CtrlChans[id] <- Message{Type: ErrMsg, Message: []byte(msg.Error())}
			}
		}
	}
}

func (s *ServerManager) handleErrors(id int) {
	for {
		select {
		case msg := <-s.ErrChans[id]:
			if msg != nil {
				s.MsgChans[id] <- Message{Type: ErrMsg, Message: []byte(msg.Error())}
			}
		}
	}
}

//Register Returns the LogMessage channel of an associated server
func (s *ServerManager) Register(id int) chan Message {
	if _, ok := s.MsgChans[id]; ok {
		return s.MsgChans[id]
	}
	return nil
}

//StopServer Stop a server specified by IP
func (s *ServerManager) StopServer(id int) error {
	//stop instance of server based on ID
	if _, ok := s.CtrlChans[id]; ok {
		s.CtrlChans[id] <- CtrlChanMsg{Type: StopMsg}
		close(s.CtrlChans[id])
	}
	return nil
}

//Finds an unused ID
func (s *ServerManager) getID() (id int) {
	rand.Seed(time.Now().Unix() * rand.Int63())
	id = rand.Int()

	if _, ok := s.CtrlChans[id]; !ok {
		return id
	}
	return s.getID()
}
