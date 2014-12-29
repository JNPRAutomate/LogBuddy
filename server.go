package logbuddy

type Server interface {
    Listen()
    Close()
}
