package logbuddy

import (
    "net"
    "bufio"
    "log"
)

//TCPServe
type TCPServer struct {
    Type string // tcp,tcp4,tcp6
    IP string
    Port int
    listenAddr *net.TCPAddr
    ctrlChan chan string
    webChan chan WebChanMsg
    sock *net.TCPListener
}

func NewTCPServer(ctrlChan chan string, webChan chan WebChanMsg, net string, ip string, port int) *TCPServer {
    s := &TCPServer{ctrlChan: ctrlChan, webChan: webChan, Type: net, IP: ip, Port: port}
    s.setListener()
    return s
}

func (s *TCPServer) Listen() {
    var err error = nil

    s.sock, err = net.ListenTCP(s.Type,s.listenAddr)
    if err != nil {
        panic(err)
    }
    for {
        conn, err := s.sock.AcceptTCP()
        if err != nil {
            log.Println("TCP ERROR",err)
        }
        go s.handleClient(conn)
    }
}

func (s *TCPServer) handleClient(conn *net.TCPConn) {
    scanner := bufio.NewScanner(conn)
    for {
        if ok := scanner.Scan(); !ok {
            break
        }
        log.Println("TCP",scanner.Text())
    }
}


func (s *TCPServer) Close() {
    if s.sock == nil {
        return
    }
    s.sock.Close()
}

func (s *TCPServer) setListener() error {
    s.listenAddr = &net.TCPAddr{IP:net.ParseIP(s.IP),Port:s.Port}
    return nil
}
