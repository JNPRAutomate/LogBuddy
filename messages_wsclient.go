package logbuddy

import (
	"errors"
	"fmt"
)

//WSClientMessage A universal message to send to web socket clients
type WSClientMessage struct {
	Type int    `json:"type"` //Type of Message
	ID   int    `json:"id"`   //ID associated with message
	Data []byte `json:"data"` //Payload of Message
}

//AddDataPayload adds a []byte payload to Data
func (wscm *WSClientMessage) AddDataPayload(data []byte) error {
	if len(data) > 0 {
		wscm.Data = data
		return nil
	}
	return errors.New("Unable to add payload")
}

//MarshalJSON returns a JSON version of WSClientMessage
func (wscm *WSClientMessage) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"type\":%d, \"id\":%d,\"data\":%s}", wscm.Type, wscm.ID, wscm.Data)), nil

}
