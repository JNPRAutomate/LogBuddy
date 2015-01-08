package logbuddy

//Server interface defining the behaviors of a server
type Server interface {
	Listen()
	close()
	setListener() error
}

//ServerConfig Specifies the network configuration of the server
type ServerConfig struct {
	Type string `json:"type"` // udp,udp4,udp6,tcp,tcp4,tcp6
	IP   string `json:"ip"`   // IP Address in string format
	Port int    `json:"port"` // Port to listen in, < 1024 requires root permission
}
