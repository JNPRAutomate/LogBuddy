package logbuddy

import (
	"log"
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

	msgChan := make(chan Message)
	ctrlChan := make(chan CtrlChanMsg)

	go func(msgChan chan Message) {
		for {
			select {
			case msg := <-msgChan:
				log.Println(string(msg.Message))
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

func SendTCPMessage(dstIP string, dstPort int, netType string, itter int, counter int, t *testing.T) {
	log.Println("Starting TCP connection:", counter)
	testConn, err := net.DialTCP(netType, nil, &net.TCPAddr{IP: net.ParseIP(dstIP), Port: dstPort})
	if err != nil {
		log.Println("TCP Client Error: ", err)
		t.Fail()
	} else {
		var i int
		for i = 0; i < 10; i++ {
			msg := "Hello " + strconv.Itoa(counter) + "\n"
			_, err := testConn.Write([]byte(msg))
			if err != nil {
				log.Println("TCP Write Error: ", err)
				t.Fail()
			}
		}
	}
	log.Println("Stopping TCP connection:", counter)
	testConn.Close()
}
