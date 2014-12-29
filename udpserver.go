package logbuddy

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
    listenAddr *net.UDPAddr
    ctrlChan chan string
    webChan chan WebChanMsg
    sock *net.UDPConn
}

func NewUDPServer(ctrlChan chan string, webChan chan WebChanMsg, ip string, port int, net string) *UDPServer {
    s := &UDPServer{ctrlChan: ctrlChan, webChan: webChan, IP: ip, Port: port, Type: net}
    s.setListener()
    return s
}

//Listen Listen to
func (s *UDPServer) Listen() {
    var buffer []byte = make([]byte,9600)
    var err error
    s.sock, err = net.ListenUDP(s.Type,s.listenAddr)
    s.sock.SetReadBuffer(1048576)
    if err != nil {
        fmt.Println(err)
    }
    for {
        //handle each packet in a seperate go routine
        s.sock.ReadFromUDP(buffer)
        go s.handlePacket(bytes.Trim(buffer,"\x00"))
    }
}

func (s *UDPServer) handlePacket(buffer []byte) {
    log.Println("UDP",len(buffer),string(buffer))
}


func (s *UDPServer) Stop() {
    s.sock.Close()
}

func (s *UDPServer) setListener() error {
    s.listenAddr = &net.UDPAddr{IP:net.ParseIP(s.IP),Port:s.Port}
    return nil
}
