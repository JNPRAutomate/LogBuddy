package logbuddy

type MsgRouter struct {
	Handlers map[int]*MsgHandler
}

type MsgHandler struct {
	Sender    chan interface{}
	Receivers []chan interface{}
}

//Register
func (mr *MsgRouter) RegisterSender(s chan interface{}, id int) {
	mh := &MsgHandler{Sender: s}
	mr.Handlers[id] = mh
}

//Unregister
func (mr *MsgRouter) UnregisterSender(id int) {
	//Remove Sender
}

func (mr *MsgRouter) RegisterReceivers(r chan interface{}, id int) int {
	mr.Handlers[id].Receivers = append(mr.Handlers[id].Receivers, r)
	return len(mr.Handlers[id].Receivers) - 1

}

//FanOut

//FanIn
