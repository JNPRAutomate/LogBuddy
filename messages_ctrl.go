package logbuddy

import "fmt"

//CtrlChanMsg control channel message format
type CtrlChanMsg struct {
	Type    int    //Type message type
	Message []byte //Message string message to include with msg
}

//String return string format of message
func (c *CtrlChanMsg) String() string {
	return fmt.Sprintf("Type: %d Message: %s", c.Type, string(c.Message))
}

//MarshalJSON returns json []byte of CtrlChanMsg
func (c *CtrlChanMsg) MarshalJSON() ([]byte, error) {
	//Check to see if inner message is a JSON message
	if c.Message[0] == byte('\x7B') {
		return []byte(fmt.Sprintf("{\"type\": %d, \"message\":%s}", c.Type, string(c.Message))), nil
	}
	return []byte(fmt.Sprintf("{\"type\": %d, \"message\":\"%s\"}", c.Type, string(c.Message))), nil
}
