package logbuddy

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

//UDPServer A server to listen for UDP messages
type UDPServer struct {
	Config     *ServerConfig
	listenAddr *net.UDPAddr
	ctrlChan   chan CtrlChanMsg
	errChan    chan error
	msgChan    chan Message
	sock       *net.UDPConn
}

//NewUDPServer creates a new initialized UDP server
func NewUDPServer(errChan chan error, ctrlChan chan CtrlChanMsg, msgChan chan Message, config *ServerConfig) *UDPServer {
	s := &UDPServer{errChan: errChan, ctrlChan: ctrlChan, msgChan: msgChan, listenAddr: nil, Config: config}
	s.setListener()
	return s
}

//Listen Listen to
func (s *UDPServer) Listen() {
	var err error
	buffer := make([]byte, 9600)

	go func() {
		for {
			select {
			case msg := <-s.ctrlChan:
				if msg.Type == StopMsg {
					s.close()
					s.errChan <- nil
				}
			}
		}
	}()

	s.sock, err = net.ListenUDP(s.Config.Type, s.listenAddr)
	if err != nil {
		s.errChan <- err
		return
	}
	s.sock.SetReadBuffer(MaxReadBuffer)
	s.msgChan <- Message{Type: AckStartMsg, Message: []byte(fmt.Sprintf("Server started: %s %s %d", s.Config.Type, s.Config.IP, s.Config.Port))}
	for {
		//handle each packet in a seperate go routine
		_, _, err := s.sock.ReadFromUDP(buffer)
		if err != nil {
			switch err := err.(type) {
			case net.Error:
				if err.Timeout() {
					s.errChan <- err
				} else if err.Temporary() {
					s.errChan <- err
				}
			}
			s.errChan <- err
			return
		}
		go func() {
			srcIP, srcPort, err := net.SplitHostPort(s.sock.LocalAddr().String())
			if err != nil {

			}
			//dstIP, dstPort, err := net.SplitHostPort(s.sock.RemoteAddr().String())
			if err != nil {

			}
			srcPortInt, err := strconv.Atoi(srcPort)
			if err != nil {

			}
			//dstPortInt, err := strconv.Atoi(dstPort)
			if err != nil {

			}
			s.msgChan <- Message{Type: DataMsg, SrcIP: net.ParseIP(srcIP), SrcPort: srcPortInt, DstIP: nil /*net.ParseIP(dstIP)*/, DstPort: 0 /*dstPortInt*/, Network: s.sock.LocalAddr().Network(), Message: bytes.Trim(buffer, "\x00")}
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
