package logbuddy

const (
	//MaxReadBuffer Max mempory allocated to a server socket for listening
	MaxReadBuffer = 1048576

	//Control message types

	//InitMsg Initialization message
	InitMsg = 0
	//DataMsg Data message
	DataMsg = 1
	//StartMsg Start message
	StartMsg = 100
	//StopMsg Stop message
	StopMsg = 255
)
