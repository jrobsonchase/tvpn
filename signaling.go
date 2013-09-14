package tvpn

type Backend interface {
	SendMessage(Message) error
	RecvMessage() Message
}
