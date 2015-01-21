package logbuddy

import (
	"errors"
	"fmt"
	"html"
	"net"
	"strings"
)

const (
	//Control Message types

	//BadMsg represents bad messages
	BadMsg = -1
	//InitMsg Initialization Message
	InitMsg = 0
	//DataMsg Data Message
	DataMsg = 1
	//ReqMsg New request Message
	ReqMsg = 2
	//ErrMsg Error Message
	ErrMsg = 3
	//CtrlMsg Control Message
	CtrlMsg = 4
	//RestartMsg Restart Message
	RestartMsg = 5
	//StartMsg Start Message
	StartMsg = 100
	//AckStartMsg acknowledge start Message
	AckStartMsg = 101
	//StopMsg Stop Message
	StopMsg = 255
)

//LogMessage A message containing a log
type LogMessage struct {
	Message []byte `json:"Message"`  //Message
	SrcIP   net.IP `json:"srcip"`    //Src IP of Message
	SrcPort int    `json:"srcport"`  //Src Port of Message
	DstIP   net.IP `json:"destip"`   //Dst IP of Message
	DstPort int    `json:"destport"` //Dest Port of Message
	Network string `json:"network"`  //Network type of Message tcp,tcp4,tcp6,udp,udp4,udp6
}

//String returns string value of LogMessage
func (m *LogMessage) String() string {
	return fmt.Sprintf("SrcIP=\"%s\" SrcPort=\"%d\" DstIP=\"%s\" DstPort=\"%d\" Network=\"%s\" Message=\"%s\"", m.SrcIP.String(), m.SrcPort, m.DstIP.String(), m.DstPort, m.Network, string(m.Message))
}

//MarshalJSON returns json []byte of LogMessage
func (m *LogMessage) MarshalJSON() ([]byte, error) {
	//check for nil values
	message := html.EscapeString(strings.TrimRight(string(m.Message), "\n"))
	return []byte(fmt.Sprintf("{ message\":\"%s\",\"srcip\":\"%s\",\"srcport\":%d,\"dstip\":\"%s\",\"dstport\":%d,\"network\":\"%s\"}", message, m.SrcIP.String(), m.SrcPort, m.DstIP.String(), m.DstPort, m.Network)), nil
}

//ClientMessage Messages sent from the websocket client
type ClientMessage struct {
	Type         int          `json:"type"`         //Type Message type
	Channel      int          `json:"channel"`      // channel to listen on
	ServerConfig ServerConfig `json:"serverconfig"` //ServerConfig configuration of requested server
}

//WSClientMessage A universal message to send to web socket clients
type WSClientMessage struct {
	Type int    `json:"type"` //Type of Message
	Data []byte `json:"data"` //Payload of Message
}

//AddDataPayload adds a []byte payload to Data
func (wscm *WSClientMessage) AddDataPayload(data []byte) error {
	if len(data) > 0 {
		wscm.Data = data
		return nil
	}
	return errors.New("Unable to add payload")
}

//MarshalJSON returns a JSON version of WSClientMessage
func (wscm *WSClientMessage) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"type\":%d,\"data\":%s}", wscm.Type, wscm.Data)), nil

}
