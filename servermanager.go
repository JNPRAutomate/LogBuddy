package logbuddy

type ServerManager struct {
	CtrlChans map[int]chan Message
}
