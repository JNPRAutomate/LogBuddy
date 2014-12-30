package logbuddy

import (
	"bufio"
	"log"
	"net"
)

//TCPServer A server to listen for TCP messages
type TCPServer struct {
	Config     *ServerConfig
	listenAddr *net.TCPAddr
	ctrlChan   chan CtrlChanMsg
	msgChan    chan Message
	sock       *net.TCPListener
}

//NewTCPServer Creates a new initialized TCPServer
func NewTCPServer(ctrlChan chan CtrlChanMsg, msgChan chan Message, config *ServerConfig) *TCPServer {
	s := &TCPServer{ctrlChan: ctrlChan, msgChan: msgChan, Config: config}
	s.setListener()
	return s
}

//Listen Starts TCPServer ready to receive messages
// Typically run as a go routine
func (s *TCPServer) Listen() error {
	go func() {
		for {
			select {
			case msg := <-s.ctrlChan:
				if msg.Type == StopMsg {
					log.Println("STOP")
					s.close()
				}

			}
		}
	}()

	var err error

	s.sock, err = net.ListenTCP(s.Config.Type, s.listenAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := s.sock.AcceptTCP()
		if err != nil {
			log.Println("TCP ERROR", err)
		}
		go func() {
			conn.SetReadBuffer(MaxReadBuffer)
			scanner := bufio.NewScanner(conn)
			for {
				if ok := scanner.Scan(); !ok {
					break
				}
				s.msgChan <- Message{Type: DataMsg, Message: scanner.Bytes()}
			}
		}()
	}
}

//Close Close the TCP Server
func (s *TCPServer) close() {
	if s.sock == nil {
		return
	}
	s.sock.Close()
}

func (s *TCPServer) setListener() error {
	s.listenAddr = &net.TCPAddr{IP: net.ParseIP(s.Config.IP), Port: s.Config.Port}
	return nil
}
