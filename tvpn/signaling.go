package tvpn

type Backend interface {
	SendMessage(Message)
	RecvMessage() Message
}
