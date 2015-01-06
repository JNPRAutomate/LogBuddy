package logbuddy

const (
	//MaxReadBuffer Max mempory allocated to a server socket for listening
	MaxReadBuffer = 1048576

	//Control message types

	//InitMsg Initialization message
	InitMsg = 0
	//DataMsg Data message
	DataMsg = 1
	//ReqMsg New request message
	ReqMsg = 2
	//StartMsg Start message
	StartMsg = 100
	//AckStartMsg acknowledge start message
	AckStartMsg = 101
	//StopMsg Stop message
	StopMsg = 255
)
