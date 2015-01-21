package logbuddy

import (
	"net"
	"strconv"
	"testing"
	"time"
)

var UDPtestPort = 5000

//Test that UDP Server Matches Interface
func TestUDPServerInterface(t *testing.T) {
	var _ Server = (*UDPServer)(nil)
}

func TestBasicUDPListener(t *testing.T) {
	counter := 0

	msgChan := make(chan LogMessage)
	ctrlChan := make(chan CtrlChanMsg)

	go func(msgChan chan LogMessage) {
		for {
			select {
			case msg := <-msgChan:
				t.Logf("%s", string(msg.Message))
			}
		}
	}(msgChan)

	listener := &UDPServer{ctrlChan: ctrlChan, msgChan: msgChan, Config: &ServerConfig{IP: "0.0.0.0", Port: UDPtestPort, Type: "udp4"}}
	listener.setListener()
	go listener.Listen()
	time.Sleep(1 * time.Second)
	for {
		counter = counter + 1
		SendUDPMessage("127.0.0.1", UDPtestPort, "udp4", 10, counter, t)
		time.Sleep(1 * time.Second)
		if counter == 10 {
			ctrlChan <- CtrlChanMsg{Type: StopMsg}
			break
		}
	}
}

func SendUDPMessageBench(dstIP string, dstPort int, netType string, itter int, counter int, t *testing.B) {
	t.Logf("%s %d", "Starting UDP connection:", counter)
	testConn, _ := net.DialUDP(netType, nil, &net.UDPAddr{IP: net.ParseIP(dstIP), Port: dstPort})
	var i int
	for i = 0; i < itter; i++ {
		_, err := testConn.Write([]byte("Hello " + strconv.Itoa(counter)))
		if err != nil {
			t.Logf("%s %s", "UDP Write Error: ", err)
			t.Fail()
		}
	}
	t.Logf("%s %d", "Stopping UDP connection:", counter)
	testConn.Close()
}

func SendUDPMessage(dstIP string, dstPort int, netType string, itter int, counter int, t *testing.T) {
	t.Logf("%s %d", "Starting UDP connection:", counter)
	testConn, _ := net.DialUDP(netType, nil, &net.UDPAddr{IP: net.ParseIP(dstIP), Port: dstPort})
	var i int
	for i = 0; i < itter; i++ {
		_, err := testConn.Write([]byte("Hello " + strconv.Itoa(counter)))
		if err != nil {
			t.Logf("%s %s", "UDP Write Error: ", err)
			t.Fail()
		}
	}
	t.Logf("%s %d", "Stopping UDP connection:", counter)
	testConn.Close()
}
