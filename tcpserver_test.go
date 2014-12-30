package logbuddy

import (
	"log"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestBasicTCPListener(t *testing.T) {
	var counter int
	counter = 0

	msgChan := make(chan Message)

	go func(msgChan chan Message) {
		for {
			select {
			case msg := <-msgChan:
				log.Println(string(msg.Message))
			}
		}
	}(msgChan)

	listener := &TCPServer{Type: "tcp4", IP: "0.0.0.0", Port: 5000, msgChan: msgChan}
	listener.setListener()
	go listener.Listen()
	time.Sleep(1 * time.Second)
	for {
		counter = counter + 1
		log.Println("Starting TCP connection:", counter)
		testConn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5000})
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
		time.Sleep(1 * time.Second)
		if counter == 10 {
			break
		}
	}
}
