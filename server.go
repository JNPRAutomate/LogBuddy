package logbuddy

import "fmt"

//Server interface defining the behaviors of a server
type Server interface {
	Listen()
	close()
	setListener() error
}

//ServerConfig Specifies the network configuration of the server
type ServerConfig struct {
	ID   int    `json:"id"`   //ID ID of running server, assigned by the server manager
	Type string `json:"type"` // udp,udp4,udp6,tcp,tcp4,tcp6
	IP   string `json:"ip"`   // IP Address in string format
	Port int    `json:"port"` // Port to listen in, < 1024 requires root permission
}

//MarshalJSON Returns json marshalled LogMessage
func (sc *ServerConfig) MarshalJSON() ([]byte, error) {
	//check for nil values
	return []byte(fmt.Sprintf("{\"id\":%d,\"type\":\"%s\",\"ip\":\"%s\",\"port\":%d}", sc.ID, sc.Type, sc.IP, sc.Port)), nil
}
