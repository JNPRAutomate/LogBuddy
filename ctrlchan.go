package logbuddy

type CtrlChanMsg struct {
	Type    int    //LogMessage type
	Message []byte //Message string message to include with msg
}
