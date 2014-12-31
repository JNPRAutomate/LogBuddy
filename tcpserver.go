package logbuddy

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"
)

//TCPServer A server to listen for TCP messages
type TCPServer struct {
	Config     *ServerConfig
	listenAddr *net.TCPAddr
	ctrlChan   chan CtrlChanMsg
	msgChan    chan Message
	sock       *net.TCPListener
	conns      []*net.TCPConn
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
	var err error

	go func() {
		for {
			select {
			case msg := <-s.ctrlChan:
				if msg.Type == StopMsg {
					log.Println("Stopping TCP Server")
					s.close()
					return
				}
			}
		}
	}()

	s.sock, err = net.ListenTCP(s.Config.Type, s.listenAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := s.sock.AcceptTCP()
		if err != nil {
			switch err := err.(type) {
			case net.Error:
				if err.Timeout() {
					log.Println("Timeout Error")
				} else if err.Temporary() {
					log.Println("Temp Error")
				}
			}
			return err
		}
		conn.SetReadBuffer(MaxReadBuffer)
		s.conns = append(s.conns, conn)
		scanner := bufio.NewScanner(conn)
		go func() {
			for {
				if ok := scanner.Scan(); !ok {
					conn.SetReadDeadline(time.Now())
					if _, err := conn.Read(make([]byte, 1)); err == io.EOF {
						conn.Close()
						conn = nil
					} else {
						conn.SetReadDeadline(time.Time{})
					}
					return
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
	for item := range s.conns {
		if s.conns[item] != nil {
			s.conns[item].SetReadDeadline(time.Now())
			if _, err := s.conns[item].Read(make([]byte, 1)); err == io.EOF {
				s.conns[item].Close()
				s.conns[item] = nil
			} else {
				s.conns[item].SetReadDeadline(time.Time{})
			}
		}

	}
	s.sock.Close()
}

func (s *TCPServer) setListener() error {
	s.listenAddr = &net.TCPAddr{IP: net.ParseIP(s.Config.IP), Port: s.Config.Port}
	return nil
}
