package logbuddy

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

//ServerManager Manages listening servers
type ServerManager struct {
	ServerConfigs map[int]*ServerConfig
	CtrlChans     map[int]chan CtrlChanMsg
	//cCtrlChans client provided control chans
	cCtrlChans map[int]chan CtrlChanMsg
	MsgChans   map[int]chan LogMessage
	ErrChans   map[int]chan error
	msgRouter  *MsgRouter
	wg         sync.WaitGroup
}

//NewServerManager Creates a new server manager with an initalized CtrlChans map
func NewServerManager() *ServerManager {
	return &ServerManager{CtrlChans: make(map[int]chan CtrlChanMsg),
		cCtrlChans:    make(map[int]chan CtrlChanMsg),
		MsgChans:      make(map[int]chan LogMessage),
		msgRouter:     &MsgRouter{},
		ErrChans:      make(map[int]chan error),
		ServerConfigs: make(map[int]*ServerConfig)}
}

//StartServer Start a new server with a server config
func (s *ServerManager) StartServer(config *ServerConfig) (id int, err error) {
	//Get a new unused ID
	id = s.getID()
	config.ID = id
	//TODO: CHECK FOR LISTENING PORT
	//TODO: ADD CLIENT SPECIFIC CTRL CHANNEL
	if config.Type == "tcp4" || config.Type == "tcp6" || config.Type == "tcp" {
		s.MsgChans[id] = make(chan LogMessage)
		s.CtrlChans[id] = make(chan CtrlChanMsg)
		s.cCtrlChans[id] = make(chan CtrlChanMsg)
		s.ErrChans[id] = make(chan error)
		listener := &TCPServer{Config: config, msgChan: s.MsgChans[id], ctrlChan: s.CtrlChans[id], errChan: s.ErrChans[id]}
		s.ServerConfigs[id] = config
		s.handleCtrl(id)
		s.handleErrors(id)
		listener.setListener()
		s.wg.Add(1)
		go listener.Listen()
		return id, nil
	} else if config.Type == "udp4" || config.Type == "udp6" || config.Type == "udp" {
		s.MsgChans[id] = make(chan LogMessage)
		s.CtrlChans[id] = make(chan CtrlChanMsg)
		s.cCtrlChans[id] = make(chan CtrlChanMsg)
		s.ErrChans[id] = make(chan error)
		s.ServerConfigs[id] = config
		listener := &UDPServer{Config: config, msgChan: s.MsgChans[id], ctrlChan: s.CtrlChans[id], errChan: s.ErrChans[id]}
		s.handleCtrl(id)
		s.handleErrors(id)
		listener.setListener()
		s.wg.Add(1)
		go listener.Listen()
		return id, nil
	}
	return id, errors.New("Unable to create specified server")
}

func (s *ServerManager) handleCtrl(id int) {
	go func(ctrlMsg <-chan CtrlChanMsg) {
		defer s.wg.Done()
		s.wg.Add(1)
		for msg := range ctrlMsg {
			log.Println("SM Ctrl", s.CtrlChans[id], msg.String())
			if msg.Type == AckStartMsg {
				msg.Message, _ = s.ServerConfigs[id].MarshalJSON()
				s.cCtrlChans[id] <- msg
			}
		}
	}(s.CtrlChans[id])
}

func (s *ServerManager) handleErrors(id int) {
	go func(errMsg <-chan error) {
		defer s.wg.Done()
		s.wg.Add(1)
		for msg := range errMsg {
			//TODO: Identify if server has stopped
			log.Println(msg)
		}
	}(s.ErrChans[id])
}

//Register Returns the LogMessage channel of an associated server
func (s *ServerManager) Register(id int) (chan LogMessage, chan CtrlChanMsg) {
	//TODO: ADD CLIENT SPECIFIC CTRL CHANNEL
	var msgChan chan LogMessage
	var ctrlChan chan CtrlChanMsg

	if _, ok := s.MsgChans[id]; ok {
		msgChan = s.MsgChans[id]
	}

	if _, ok := s.cCtrlChans[id]; ok {
		ctrlChan = s.cCtrlChans[id]
	}
	return msgChan, ctrlChan
}

//StopServer Stop a server specified by IP
func (s *ServerManager) StopServer(id int) error {
	//stop instance of server based on ID
	if _, ok := s.CtrlChans[id]; ok {
		s.CtrlChans[id] <- CtrlChanMsg{Type: StopMsg}
		s.wg.Done()
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
