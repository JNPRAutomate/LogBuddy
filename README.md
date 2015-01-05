App Usage
===========

Pending inital release.

WebSocket API
=============

All client communication to the main server is done via WebSockets. This keeps communication cleanly running between the client and the webserver. Also the core point of the tool is to view logs in the web UI. To maximize the intake of the logs using seperate WebSocket connections allows us to consume data at a faster rate.

Const
-----

These are constants used for the JSON messages

### MSG Types

```
//Init message - For creating a server
var INIT_MSG = 0
//A message containing data
var DATA_MSG = 1
//Start message - For starting a stopped server
var START_MSG = 100
//Stop message -  For stopping a server
var STOP_MSG = 255
```

API
---

All calls are sent as a WebSocket TextMessage type messages with a JSON payload.

### Server Config

```
{
	Type: "udp" //Specify the type of server to listen "udp","udp4","udp6","tcp","tcp4","tcp6"
	IP: "0.0.0.0" //IP String to listen on, "0.0.0.0" means all interfaces
	Port 5000 //Port to listen on, listening on port < 1024 requires root
}
```

Basic JSON Format
-----------------

```
{
	Type: INIT_MSG  //Required - Defines message type
	serverConfig:
}
```

### Control Sessions
