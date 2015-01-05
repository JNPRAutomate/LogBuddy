package logbuddy

import (
	"fmt"
	"net"
)

//Message struct used to pass logging information
type Message struct {
	Type    int    `json:"type"`     //Msg type
	Message []byte `json:"message"`  //message
	SrcIP   net.IP `json:"srcip"`    //Src IP of message
	SrcPort int    `json:"srcport"`  //Src Port of message
	DstIP   net.IP `json:"destip"`   //Dst IP of message
	DstPort int    `json:"destport"` //Dest Port of message
	Network string `json:"network"`  //Network type of message tcp,tcp4,tcp6,udp,udp4,udp6
}

func (m *Message) String() string {
	return fmt.Sprintf("Type=\"%d\" SrcIP=\"%s\" SrcPort=\"%d\" DstIP=\"%s\" DstPort=\"%d\" Network=\"%s\" Message=\"%s\"", m.Type, m.SrcIP.String(), m.SrcPort, m.DstIP.String(), m.DstPort, m.Network, string(m.Message))
}

//ClientMessage Messages sent from the websocket client
type ClientMessage struct {
	Type    int //Type message type
	Channel int // channel to listen on
}
