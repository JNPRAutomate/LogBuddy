//Constants for LogMessage types
var INIT_MSG = 0;
var DATA_MSG = 1;
var REQ_MSG = 2;
var ERR_MSG = 3;
var CTRL_MSG = 4;
var RESTART_MSG = 5;
var START_MSG = 100;
var ACK_START_MSG = 101;
var STOP_MSG = 255;

/*
	JSONinc golang like object structure

	WSClientLogMessage : {
		Type: //specifys the type of payload
		Data: //payload json
	}
*/

//ClientLogMessage generates a new clientLogMessage
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
	//handle close
	this.conn.onclose = function(evt) {
		$("#status-icon").addClass("conn-status-disconn");
		$("#status-icon").removeClass("conn-status-conn");
	}
	//handle LogMessages
	this.conn.onmessage = function(evt) {
		console.log(evt.data)
		msg = JSON.parse(evt.data);
		if (msg.type === DATA_MSG) {
			$("#log-display").append($("<p>", {html: msg.data.message}));
			console.log(msg.data.message);
		} else if (msg.type === ACK_START_MSG) {
			$("#server-list").append("<li data-id=\""+msg.data.id+"\">IP:"+msg.data.ip+" Port:"+msg.data.port+" Type: "+msg.data.type+"</li>");
			//register started server
			console.log(msg.message);
		} else if (msg.type === ERR_MSG) {
			//remove server if error thrown
			console.log(msg.message);
		} else if (msg.type === RESTART_MSG) {
			$("#server-list").append("<li data-id=\""+msg.data.id+"\">IP:"+msg.data.ip+" Port:"+msg.data.port+" Type: "+msg.data.type+"</li>");
		}
	}
	//handle errors
	this.conn.onerror = function(evt) {
		document.getElementById("logData").innerHTML = evt.data;
	}
	//handle open
	this.conn.onopen = function(evt) {
		$("#status-icon").removeClass("conn-status-disconn");
		$("#status-icon").addClass("conn-status-conn");
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

//Checks to see if the server is listneing
ServerConn.prototype.CheckConn = function() {
	var self = this;
	if (self.conn.readyState === 1) {
		return true;
	} else {
		return false;
	}
}
