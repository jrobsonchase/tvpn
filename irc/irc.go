package irc

import (
	"github.com/Pursuit92/irc"
	"github.com/Pursuit92/tvpn"
)

func SetLogLevel(n int) {
	irc.SetLogLevel(n)
}

type IRCBackend struct {
	Conn        *irc.Conn
	MsgExpector irc.Expector
	Messages    chan irc.Command
	Status      chan irc.Command
	Convos      chan map[string]irc.ExpectChan
}

func Connect(host, nick, group string) (*IRCBackend, error) {
	conn, err := irc.DialIRC(host, []string{nick}, nick, nick)
	if err != nil {
		return nil, err
	}
	_, err = conn.Register()
	if err != nil {
		return nil, err
	}

	chann, err := conn.Join(group)
	if err != nil {
		return nil, err
	}

	joinpart, _ := irc.Expect(chann, irc.Command{"", "(JOIN)|(PART)", []string{}})
	quit, _ := irc.Expect(conn, irc.Command{"", "QUIT", []string{}})
	status := make(chan irc.Command)
	// Combine joinpart and quit into one channel
	go func() {
		for {
			select {
			case msg := <-joinpart.Chan:
				status <- msg
			case msg := <-quit.Chan:
				status <- msg
			}
		}
	}()

	msgs, err := irc.Expect(conn, irc.Command{"", "PRIVMSG", []string{nick}})
	if err != nil {
		return nil, err
	}

	users := chann.GetUsers()
	go makeJoin(users, status)

	return &IRCBackend{Conn: conn, Messages: msgs.Chan, Status: status}, nil
}

func makeJoin(users map[string]irc.IRCUser, status chan<- irc.Command) {
	for _, v := range users {
		status <- irc.Command{v.String(), irc.Join, []string{}}
	}
}

func (b IRCBackend) RecvMessage() tvpn.Message {
	for {
		select {
		case input := <-b.Messages:
			ircMsg := input.Message()
			msg, err := tvpn.ParseMessage(input.Params[len(input.Params)-1])
			if err == nil {
				msg.From = ircMsg.Nick
				return *msg
			}
		case input := <-b.Status:
			switch input.Command {
			case "QUIT", "PART":
				return tvpn.Message{From: input.Message().Nick, Type: tvpn.Quit}
			case "JOIN":
				if input.Message().Nick != b.Conn.Nick {
					return tvpn.Message{From: input.Message().Nick, Type: tvpn.Join}
				}
			}
		}

	}
}

func (b IRCBackend) SendMessage(mes tvpn.Message) error {
	return b.Conn.Send(irc.Command{b.Conn.Nick, irc.Privmsg, []string{mes.To, mes.String()}})
}
