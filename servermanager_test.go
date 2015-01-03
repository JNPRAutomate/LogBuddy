package logbuddy

import (
	"log"
	"testing"
	"time"
)

func TestServerManager(t *testing.T) {
	ServerManagerItter := 10
	ServerManagerTestPort := 5001

	counter := 0
	//create new server manager
	sm := &ServerManager{CtrlChans: make(map[int]chan CtrlChanMsg), MsgChans: make(map[int]chan Message), ServerConfigs: make(map[int]*ServerConfig)}
	//start tcp server
	tcpServerID, err := sm.StartServer(&ServerConfig{IP: "0.0.0.0", Port: ServerManagerTestPort, Type: "tcp4"})
	log.Println("Starting TCP Server ID", tcpServerID)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	time.Sleep(1 * time.Second)

	//start udp server
	udpServerID, err := sm.StartServer(&ServerConfig{IP: "0.0.0.0", Port: ServerManagerTestPort, Type: "udp4"})
	log.Println("Starting UDP Server ID", udpServerID)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	time.Sleep(1 * time.Second)

	//send messages to each server
	for {
		counter = counter + 1
		SendTCPMessage("127.0.0.1", ServerManagerTestPort, "tcp4", ServerManagerItter, counter, t)
		if counter == ServerManagerItter {
			log.Println("Stopping TCP Server ID", tcpServerID)
			sm.StopServer(tcpServerID)
			time.Sleep(5 * time.Second)
			break
		}
	}

	//reset counter
	counter = 0
	for {
		counter = counter + 1
		SendUDPMessage("127.0.0.1", ServerManagerTestPort, "udp4", ServerManagerItter, counter, t)
		if counter == ServerManagerItter {
			log.Println("Stopping UDP Server ID", udpServerID)
			sm.StopServer(udpServerID)
			time.Sleep(5 * time.Second)
			break
		}
	}
}

func BenchmarkServerManager(b *testing.B) {
	ServerManagerItter := 1000
	ServerManagerTestPort := 5001

	counter := 0
	//create new server manager
	sm := &ServerManager{CtrlChans: make(map[int]chan CtrlChanMsg), MsgChans: make(map[int]chan Message)}
	//start tcp server
	tcpServerID, err := sm.StartServer(&ServerConfig{IP: "0.0.0.0", Port: ServerManagerTestPort, Type: "tcp4"})
	log.Println("Starting TCP Server ID", tcpServerID)
	if err != nil {
		log.Println(err)
		b.Fail()
	}
	time.Sleep(1 * time.Second)

	//start udp server
	udpServerID, err := sm.StartServer(&ServerConfig{IP: "0.0.0.0", Port: ServerManagerTestPort, Type: "udp4"})
	log.Println("Starting UDP Server ID", udpServerID)
	if err != nil {
		log.Println(err)
		b.Fail()
	}
	time.Sleep(1 * time.Second)

	//send messages to each server
	for {
		counter = counter + 1
		SendTCPMessage("127.0.0.1", ServerManagerTestPort, "tcp4", ServerManagerItter, counter, nil)
		if counter == ServerManagerItter {
			log.Println("Stopping TCP Server ID", tcpServerID)
			sm.StopServer(tcpServerID)
			time.Sleep(5 * time.Second)
			break
		}
	}

	//reset counter
	counter = 0
	for {
		counter = counter + 1
		SendUDPMessage("127.0.0.1", ServerManagerTestPort, "udp4", ServerManagerItter, counter, nil)
		if counter == ServerManagerItter {
			log.Println("Stopping UDP Server ID", udpServerID)
			sm.StopServer(udpServerID)
			time.Sleep(5 * time.Second)
			break
		}
	}
}
