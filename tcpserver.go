package logbuddy

import (
	"bufio"
	"log"
	"net"
)

//TCPServer A server to listen for TCP messages
type TCPServer struct {
	Type       string // tcp,tcp4,tcp6
	IP         string
	Port       int
	listenAddr *net.TCPAddr
	ctrlChan   chan string
	msgChan    chan Message
	sock       *net.TCPListener
}

//NewTCPServer Creates a new initialized TCPServer
func NewTCPServer(ctrlChan chan string, msgChan chan Message, net string, ip string, port int) *TCPServer {
	s := &TCPServer{ctrlChan: ctrlChan, msgChan: msgChan, Type: net, IP: ip, Port: port}
	s.setListener()
	return s
}

//Listen Starts TCPServer ready to receive messages
// Typically run as a go routine
func (s *TCPServer) Listen() error {
	var err error

	s.sock, err = net.ListenTCP(s.Type, s.listenAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := s.sock.AcceptTCP()
		if err != nil {
			log.Println("TCP ERROR", err)
		}
		go s.handleClient(conn)
	}
}

func (s *TCPServer) handleClient(conn *net.TCPConn) {
	conn.SetReadBuffer(MaxReadBuffer)
	scanner := bufio.NewScanner(conn)
	for {
		if ok := scanner.Scan(); !ok {
			break
		}
		s.msgChan <- Message{Type: DataMsg, Message: scanner.Bytes()}
	}
}

//Close Close the TCP Server
func (s *TCPServer) Close() {
	if s.sock == nil {
		return
	}
	s.sock.Close()
}

func (s *TCPServer) setListener() error {
	s.listenAddr = &net.TCPAddr{IP: net.ParseIP(s.IP), Port: s.Port}
	return nil
}
