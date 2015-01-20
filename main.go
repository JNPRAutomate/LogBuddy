package logbuddy

//go:generate go-bindata -pkg "logbuddy" -o logbuddy_bindata.go static/...

const (
	//MaxReadBuffer Max mempory allocated to a server socket for listening
	MaxReadBuffer = 1048576
)
