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
	webChan    chan WebChanMsg
	sock       *net.TCPListener
}

//NewTCPServer Creates a new initialized TCPServer
func NewTCPServer(ctrlChan chan string, webChan chan WebChanMsg, net string, ip string, port int) *TCPServer {
	s := &TCPServer{ctrlChan: ctrlChan, webChan: webChan, Type: net, IP: ip, Port: port}
	s.setListener()
	return s
}

//Listen Starts TCPServer ready to receive messages
// Typically run as a go routine
func (s *TCPServer) Listen() {
	var err error

	s.sock, err = net.ListenTCP(s.Type, s.listenAddr)
	if err != nil {
		panic(err)
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
		log.Println("TCP", scanner.Text())
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
