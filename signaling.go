package tvpn

type SigBackend interface {
	SendMessage(Message) error
	RecvMessage() Message
}
