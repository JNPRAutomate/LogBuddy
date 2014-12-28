package servers

import (
    "net"
    "log"
    "fmt"
    "bytes"
)

//UDPServer
type UDPServer struct {
    Type string // udp,udp4,udp6
    IP string
    Port int
    ListenAddr *net.UDPAddr
    ctrlChan chan string
    webChan chan WebChanMsg
    sock *net.UDPConn
}

func NewUDPServer(ctrlChan chan string, webChan chan string, ip string, port int, type string) *UDPServer {
    return &UDPServer{ctrlChan: ctrlChan, webChan: webChan, IP: ip, Port: port, Type: string}
}

//Listen Listen to
func (s *UDPServer) Listen() {
    var buffer []byte = make([]byte,9600)
    var err error
    s.sock, err = net.ListenUDP(s.Type,s.ListenAddr)
    s.sock.SetReadBuffer(1048576)
    if err != nil {
        fmt.Println(err)
    }
    for {
        log.Println("FOO")
        //handle each packet in a seperate go routine
        readBytes, _, _ := s.sock.ReadFromUDP(buffer)
        go func() {
            log.Println(readBytes,bytes.Trim(buffer,"\x00"))
        }()
    }
}

func (s *UDPServer) Stop() {
    s.sock.Close()
}
