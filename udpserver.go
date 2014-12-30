package logbuddy

import (
	"bytes"
	"log"
	"net"
)

//UDPServer A server to listen for UDP messages
type UDPServer struct {
	Type       string // udp,udp4,udp6
	IP         string
	Port       int
	listenAddr *net.UDPAddr
	ctrlChan   chan CtrlChanMsg
	msgChan    chan Message
	sock       *net.UDPConn
}

//NewUDPServer creates a new initialized UDP server
func NewUDPServer(ctrlChan chan CtrlChanMsg, msgChan chan Message, ip string, port int, net string) *UDPServer {
	s := &UDPServer{ctrlChan: ctrlChan, msgChan: msgChan, IP: ip, Port: port, Type: net, listenAddr: nil}
	s.setListener()
	return s
}

//Listen Listen to
func (s *UDPServer) Listen() error {

	go func(s *UDPServer) {
		for {
			select {
			case msg := <-s.ctrlChan:
				if msg.Type == StopMsg {
					log.Println("STOP")
					s.stop()
				}

			}
		}
	}(s)

	buffer := make([]byte, 9600)
	var err error
	s.sock, err = net.ListenUDP(s.Type, s.listenAddr)
	s.sock.SetReadBuffer(MaxReadBuffer)
	if err != nil {
		return err
	}
	for {
		//handle each packet in a seperate go routine
		s.sock.ReadFromUDP(buffer)
		go s.handlePacket(bytes.Trim(buffer, "\x00"))
	}
}

func (s *UDPServer) handlePacket(buffer []byte) {
	s.msgChan <- Message{Type: DataMsg, Message: buffer}
}

//Stop stops the UDP server from listening
func (s *UDPServer) stop() {
	s.sock.Close()
}

//setListener creates the UDPAddr for the
func (s *UDPServer) setListener() error {
	s.listenAddr = &net.UDPAddr{IP: net.ParseIP(s.IP), Port: s.Port}
	return nil
}
