package logbuddy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

//TCPServer A server to listen for TCP messages
type TCPServer struct {
	Config     *ServerConfig
	listenAddr *net.TCPAddr
	ctrlChan   chan CtrlChanMsg
	msgChan    chan Message
	errChan    chan error
	sock       *net.TCPListener
	conns      []*net.TCPConn
}

//NewTCPServer Creates a new initialized TCPServer
func NewTCPServer(errChan chan error, ctrlChan chan CtrlChanMsg, msgChan chan Message, config *ServerConfig) *TCPServer {
	s := &TCPServer{errChan: errChan, ctrlChan: ctrlChan, msgChan: msgChan, Config: config}
	s.setListener()
	return s
}

//Listen Starts TCPServer ready to receive messages
// Typically run as a go routine
func (s *TCPServer) Listen() {
	var err error

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

	s.sock, err = net.ListenTCP(s.Config.Type, s.listenAddr)
	if err != nil {
		s.errChan <- err
		return
	}
	s.msgChan <- Message{Type: AckStartMsg, Message: []byte(fmt.Sprintf("Server started: %s %s %d", s.Config.Type, s.Config.IP, s.Config.Port))}
	for {
		conn, err := s.sock.AcceptTCP()
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
				srcIP, srcPort, err := net.SplitHostPort(conn.RemoteAddr().String())
				if err != nil {
					break
				}
				dstIP, dstPort, err := net.SplitHostPort(conn.LocalAddr().String())
				if err != nil {
					break
				}
				srcPortInt, err := strconv.Atoi(srcPort)
				if err != nil {
					break
				}
				dstPortInt, err := strconv.Atoi(dstPort)
				if err != nil {
					break
				}
				s.msgChan <- Message{Type: DataMsg, SrcIP: net.ParseIP(srcIP), SrcPort: srcPortInt, DstIP: net.ParseIP(dstIP), DstPort: dstPortInt, Network: conn.LocalAddr().Network(), Message: scanner.Bytes()}
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
