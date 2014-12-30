package logbuddy

import (
	"bytes"
	"log"
	"net"
)

var UDPTestPort = 5000

//UDPServer A server to listen for UDP messages
type UDPServer struct {
	Config     *ServerConfig
	listenAddr *net.UDPAddr
	ctrlChan   chan CtrlChanMsg
	msgChan    chan Message
	sock       *net.UDPConn
}

//NewUDPServer creates a new initialized UDP server
func NewUDPServer(ctrlChan chan CtrlChanMsg, msgChan chan Message, config *ServerConfig) *UDPServer {
	s := &UDPServer{ctrlChan: ctrlChan, msgChan: msgChan, listenAddr: nil, Config: config}
	s.setListener()
	return s
}

//Listen Listen to
func (s *UDPServer) Listen() error {

	go func() {
		for {
			select {
			case msg := <-s.ctrlChan:
				if msg.Type == StopMsg {
					log.Println("Stopping UDP Server")
					s.close()
					return
				}

			}
		}
	}()

	buffer := make([]byte, 9600)
	var err error
	s.sock, err = net.ListenUDP(s.Config.Type, s.listenAddr)
	s.sock.SetReadBuffer(MaxReadBuffer)
	if err != nil {
		return err
	}
	for {
		//handle each packet in a seperate go routine
		_, _, err := s.sock.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
		}
		go func() {
			s.msgChan <- Message{Type: DataMsg, Message: bytes.Trim(buffer, "\x00")}
		}()
	}
}

//Stop stops the UDP server from listening
func (s *UDPServer) close() {
	s.sock.Close()
}

//setListener creates the UDPAddr for the
func (s *UDPServer) setListener() error {
	s.listenAddr = &net.UDPAddr{IP: net.ParseIP(s.Config.IP), Port: s.Config.Port}
	return nil
}
