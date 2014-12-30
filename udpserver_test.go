package logbuddy

import (
	"log"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestBasicUDPListener(t *testing.T) {
	var counter int
	counter = 0

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

	listener := &UDPServer{ctrlChan: ctrlChan, msgChan: msgChan, Type: "udp4", IP: "0.0.0.0", Port: 5000, listenAddr: &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 5000}}
	go listener.Listen()
	time.Sleep(1 * time.Second)
	for {
		counter = counter + 1
		testConn, _ := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5000})
		var i int
		for i = 0; i < 10; i++ {
			_, err := testConn.Write([]byte("Hello " + strconv.Itoa(counter)))
			if err != nil {
				log.Println("UDP Write Error: ", err)
				t.Fail()
			}
		}
		testConn.Close()
		time.Sleep(1 * time.Second)
		if counter == 10 {
			ctrlChan <- CtrlChanMsg{Type: StopMsg}
			break
		}
	}
}
