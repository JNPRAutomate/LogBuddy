package logbuddy

const (
	MaxReadBuffer = 1048576 //MaxReadBuffer Max mempory allocated to a server socket for listening
	//Control Channel Message
	InitMsg  = 0
	DataMsg  = 1
	StartMsg = 100
	StopMsg  = 255
)
