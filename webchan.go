package logbuddy

import (
	"net"
)

type WebChanMsg struct {
	Type     int    //Msg type
	Message  []byte //message
	SrcIP    net.IP //Src IP of message
	SrcPort  int    //Src Port of message
	DestIP   net.IP //Dst IP of message
	DestPort int    //Dest Port of message
	Network  string //Network type of message
}
