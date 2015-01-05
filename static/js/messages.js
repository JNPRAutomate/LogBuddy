//Constants for message types
var INIT_MSG = 0
var DATA_MSG = 1
var START_MSG = 100
var STOP_MSG = 255

//ClientMessage the structure for a client message
var ClientMessage = {
	Type: -1,
	Channel: 0
}

var NewClientMessage = function(type, channel) {
	cm = ClientMessage = {Type: type, Channel: channel}
	return cm
}
