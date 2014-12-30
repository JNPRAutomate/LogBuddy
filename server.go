package logbuddy

//Server interface defining the behaviors of a server
type Server interface {
	Listen() error
	close()
	setListener() error
}

//ServerConfig Specifies the network configuration of the server
type ServerConfig struct {
	Type string // udp,udp4,udp6,tcp,tcp4,tcp6
	IP   string // IP Address in string format
	Port int
}
