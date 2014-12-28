package servers

type Server interface {
    Listen()
    Stop()
}
