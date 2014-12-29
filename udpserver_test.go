package logbuddy

import (
    "testing"
    "net"
    "log"
    "time"
    "strconv"
)

func TestBasicUDPListener(t *testing.T) {
    var counter int
    counter = 0
    listener := &UDPServer{Type:"udp4",IP:"0.0.0.0",Port: 5000}
    listener.setListener()
    go listener.Listen()
    time.Sleep(1 * time.Second)
    for {
        counter = counter + 1
        testConn, _ := net.DialUDP("udp4",nil,&net.UDPAddr{IP:net.ParseIP("127.0.0.1"),Port:5000})
        var i int
        for i = 0; i < 10; i++ {
            _, err := testConn.Write([]byte("Hello " + strconv.Itoa(counter)))
            if err != nil {
                log.Println("UDP Write Error: ",err)
                t.Fail()
            }
        }
        testConn.Close()
        time.Sleep(1 * time.Second)
        if counter == 10 {
            break
        }
    }
}
