package logbuddy

const (
    InitMsg = 0
    DataMsg = 1
    StartMsg = 100
    StopMsg = 255
)

type WebChanMsg struct {
    Type int
    Message string //message
}
