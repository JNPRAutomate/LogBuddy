package logbuddy

import (
	"net"
	"strconv"
	"testing"
	"time"
)

var TCPtestPort = 5000

//Test that TCP Server Matches Interface
func TestTCPServerInterface(t *testing.T) {
	var _ Server = (*TCPServer)(nil)
}

func TestBasicTCPListener(t *testing.T) {
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

	listener := &TCPServer{Config: &ServerConfig{IP: "0.0.0.0", Port: TCPtestPort, Type: "tcp4"}, msgChan: msgChan, ctrlChan: ctrlChan}
	listener.setListener()
	go listener.Listen()
	time.Sleep(1 * time.Second)
	for {
		counter = counter + 1
		SendTCPMessage("127.0.0.1", TCPtestPort, "tcp4", 10, counter, t)
		time.Sleep(1 * time.Second)
		if counter == 10 {
			ctrlChan <- CtrlChanMsg{Type: StopMsg}
			break
		}
	}
}

func SendTCPMessageBench(dstIP string, dstPort int, netType string, itter int, counter int, t *testing.B) {
	t.Logf("%s %d", "Starting TCP connection:", counter)
	testConn, err := net.DialTCP(netType, nil, &net.TCPAddr{IP: net.ParseIP(dstIP), Port: dstPort})
	if err != nil {
		t.Logf("%s %s", "TCP Client Error: ", err)
		t.Fail()
	} else {
		var i int
		for i = 0; i < 10; i++ {
			msg := "Hello " + strconv.Itoa(counter) + "\n"
			_, err := testConn.Write([]byte(msg))
			if err != nil {
				t.Logf("%s %s", "TCP Write Error: ", err)
				t.Fail()
			}
		}
	}
	t.Logf("%s %d", "Stopping TCP connection:", counter)
	testConn.Close()
}

func SendTCPMessage(dstIP string, dstPort int, netType string, itter int, counter int, t *testing.T) {
	t.Logf("%s %d", "Starting TCP connection:", counter)
	testConn, err := net.DialTCP(netType, nil, &net.TCPAddr{IP: net.ParseIP(dstIP), Port: dstPort})
	if err != nil {
		t.Logf("%s %s", "TCP Client Error: ", err)
		t.Fail()
	} else {
		var i int
		for i = 0; i < itter; i++ {
			msg := "Hello " + strconv.Itoa(counter) + "\n"
			_, err := testConn.Write([]byte(msg))
			if err != nil {
				t.Logf("%s %s", "TCP Write Error: ", err)
				t.Fail()
			}
		}
	}
	t.Logf("%s %d", "Stopping TCP connection:", counter)
	testConn.Close()
}
