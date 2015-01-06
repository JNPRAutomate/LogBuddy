//Constants for message types
var INIT_MSG = 0;
var DATA_MSG = 1;
var REQ_MSG = 2;
var START_MSG = 100;
var ACK_START_MSG = 101;
var STOP_MSG = 255;

/*
	JSONinc golang like object structure

	ClientMessage : {
		Type:
		Channel:
		ServerConfig:
	}
*/

//ClientMessage generates a new clientmessage
var NewClientMsg = function(type, channel, serverconfig) {
	if (channel == null || channel == undefined) {
		channel = 0
	}
	return JSON.stringify({"type": type, "channel": channel, "serverconfig": serverconfig});
}

var NewServer = function(type, dstip, dstport) {
	return {"type":type,"ip":dstip,"port":dstport}
}

var ServerConn = function(){
	this.conn = new WebSocket("ws://" + document.location.host + "/logs");
	//handle open
	this.conn.onopen = function(evt) {
		document.getElementById("logData").textContent = 'Connection Open';
	}
	//handle close
	this.conn.onclose = function(evt) {
		document.getElementById("logData").textContent = 'Connection Closed';
	}
	//handle messages
	this.conn.onmessage = function(evt) {
		msg = JSON.parse(evt.data);
		if (msg.Type === DATA_MSG) {
			document.getElementById("logData").textContent = msg.message;
		} else if (msg.Type === ACK_START_MSG) {
			//register started server
		}
	}
}

//Starts a new listening server
ServerConn.prototype.StartServer = function(type, dstip, dstport){
	this.conn.send(NewClientMsg(START_MSG,0,NewServer(type,dstip,dstport)));
}

//Stops an existing server
ServerConn.prototype.StopServer = function(channel) {
	this.conn.send(NewClientMsg(STOP_MSG,channel));
}

//Requests an existing server
ServerConn.prototype.ReqChannel = function(channel) {
	this.conn.send(NewClientMsg(REQ_MSG,channel));
}
