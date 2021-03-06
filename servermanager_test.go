package logbuddy

import (
	"testing"
	"time"
)

func TestServerManager(t *testing.T) {
	ServerManagerItter := 10
	ServerManagerTestPort := 5001

	counter := 0
	//create new server manager
	sm := &ServerManager{CtrlChans: make(map[int]chan CtrlChanMsg), MsgChans: make(map[int]chan LogMessage), ServerConfigs: make(map[int]*ServerConfig), ErrChans: make(map[int]chan error)}
	//start tcp server
	tcpServerID, err := sm.StartServer(&ServerConfig{IP: "0.0.0.0", Port: ServerManagerTestPort, Type: "tcp4"})
	t.Logf("%s %d", "Starting TCP Server ID", tcpServerID)
	if err != nil {
		t.Logf("%s", err)
		t.Fail()
	}
	time.Sleep(1 * time.Second)

	//start udp server
	udpServerID, err := sm.StartServer(&ServerConfig{IP: "0.0.0.0", Port: ServerManagerTestPort, Type: "udp4"})
	t.Logf("%s %d", "Starting UDP Server ID", udpServerID)
	if err != nil {
		t.Logf("%s", err)
		t.Fail()
	}
	time.Sleep(1 * time.Second)

	//send Messages to each server
	for {
		counter = counter + 1
		SendTCPMessage("127.0.0.1", ServerManagerTestPort, "tcp4", ServerManagerItter, counter, t)
		if counter == ServerManagerItter {
			t.Logf("%s %d", "Stopping TCP Server ID", tcpServerID)
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
			t.Logf("%s %d", "Stopping UDP Server ID", udpServerID)
			sm.StopServer(udpServerID)
			time.Sleep(5 * time.Second)
			break
		}
	}
}

func BenchmarkServerManagerTCP(b *testing.B) {
	ServerManagerTestPort := 5001

	counter := 0
	//create new server manager
	sm := &ServerManager{CtrlChans: make(map[int]chan CtrlChanMsg), MsgChans: make(map[int]chan LogMessage), ServerConfigs: make(map[int]*ServerConfig), ErrChans: make(map[int]chan error)}
	//start tcp server
	tcpServerID, err := sm.StartServer(&ServerConfig{IP: "0.0.0.0", Port: ServerManagerTestPort, Type: "tcp4"})
	b.Logf("%s %d", "Starting TCP Server ID", tcpServerID)
	if err != nil {
		b.Logf("%s", err)
		b.Fail()
	}
	time.Sleep(1 * time.Second)
	b.ResetTimer()
	//send Messages to each server
	SendTCPMessageBench("127.0.0.1", ServerManagerTestPort, "tcp4", b.N, counter, b)
	b.Logf("%s %d", "Stopping TCP Server ID", tcpServerID)
	b.Log("Packets Sent:", counter)
	sm.StopServer(tcpServerID)
}
