package logbuddy

import (
	"net"
)

type WebChanMsg struct {
	Type     int    //Msg type
	LogMessage  []byte //LogMessage
	SrcIP    net.IP //Src IP of LogMessage
	SrcPort  int    //Src Port of LogMessage
	DestIP   net.IP //Dst IP of LogMessage
	DestPort int    //Dest Port of LogMessage
	Network  string //Network type of LogMessage
}
