package servers

import (
    "testing"
    "net"
    "log"
    "time"
)

func TestBasicUDPListener(t *testing.T) {
    var counter int
    counter = 0
    listener := &UDPServer{Port: 5000, Type:"udp4", ListenAddr: &net.UDPAddr{IP:net.ParseIP("0.0.0.0"),Port:5000}}
    go listener.Listen()
    for {
        testConn, _ := net.DialUDP("udp4",nil,&net.UDPAddr{IP:net.ParseIP("127.0.0.1"),Port:5000})
        counter = counter + 1
        var i int
        for i = 0; i < 1000; i++ {
          testConn.Write([]byte("Hello " + string(counter)))
        }
        bytesWritten, err := testConn.Write([]byte("Hello " + string(counter)))
        if err != nil {
            log.Println(err)
            t.Fail()
        }
        testConn.Close()
        log.Println("Bytes Written: ",bytesWritten)
        time.Sleep(1 * time.Second)
        if counter == 10 {
            //break
        }
    }
}
