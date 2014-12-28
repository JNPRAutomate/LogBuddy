package servers

import (
    "net"
    "fmt"
)

//TCPServe
type TCPServer struct {
    Type string // tcp,tcp4,tcp6
    ListenAddr *net.TCPAddr
    ctrlChan chan string
    webChan chan WebChanMsg
    sock *net.TCPListener
}

func NewTCPServer(ctrlChan chan string, webChan chan string) *TCPServer {
    return &TCPServer{ctrlChan: ctrlChan}
}

func (s *TCPServer) Listen() {
    var err error
    s.sock, err = net.ListenTCP(s.Type,s.ListenAddr)
    if err != nil {
        fmt.Println(err)
    }
    for {
        conn, err = s.sock.Accept()
        if err != nil {
            fmt.Println(err)
        }
        go s.handleClient(conn)
    }
}

func (s *TCPServer) handleClient(conn net.Conn) {
    var buffer []byte = make([]byte,9600)
    for {
        msgLen, err := conn.Read(buffer)
        if err != nil {
            fmt.Println(err)
        }
    }
}


func (s *TCPServer) Stop() {
    s.sock.Close()
}
