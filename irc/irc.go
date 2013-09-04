package irc

import (
	"strconv"
	"strings"
	"fmt"
	"irc"
	"regexp"
	"tvpn/tvpn"
)

type IRCBackend struct {
	Conn *irc.Conn
	MsgExpector irc.Expector
	Messages chan irc.Command
	Joins chan irc.Command
	Convos chan map[string]irc.ExpectChan
}

func Connect(host, nick, group string) (*IRCBackend, error) {
	conn,err := irc.DialIRC(host, []string{nick}, nick, nick)
	if err != nil {
		return nil, err
	}
	_,err = conn.Register()
	if err != nil {
		return nil, err
	}

	chann,err := conn.Join(group)
	if err != nil {
		return nil, err
	}

	joinpart,_ := irc.Expect(chann, irc.Command{"","(JOIN)|(PART)",[]string{}})
	quit,_ := irc.Expect(conn, irc.Command{"","QUIT",[]string{}})
	status := make(chan irc.Command)
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

	msgs,err := irc.Expect(conn, irc.Command{"","PRIVMSG",[]string{nick}})
	if err != nil {
		return nil, err
	}

	convos := make(map[string]irc.ExpectChan)
	users := chann.GetUsers()
	for _,v := range users {
		fmt.Printf("Nick: %s\n",v.Nick)
		fmt.Printf("Name: %s\n",v.Name)
	}


	//msgExpector := irc.MakeExpector(msgs.Chan)
	convoChan := make(chan map[string]irc.ExpectChan,1)
	convoChan <- convos

	return &IRCBackend{Conn: conn,Messages: msgs.Chan,Joins: status,Convos: convoChan},nil
}

func parseMultiPart(message irc.Command) (string,int) {
	var remaining int
	var content string

	mpartParser := regexp.MustCompile(`^MPART (?P<remaining>[0-9]+) (?P<content>.*)$`)

	body := message.Params[len(message.Params)-1]
	if mpartParser.MatchString(body) {
		content = mpartParser.ReplaceAllString(body,"${content}")
		remaining,_ = strconv.Atoi(mpartParser.ReplaceAllString(body,"${remaining}"))
	} else {
		content = body
		remaining = 0
	}
	return content,remaining
}

// TODO - this doesn't work with multiple senders. Also doesn't supply sender
func (b IRCBackend) RecvMessage() *tvpn.Message {
	for {
		parts := make([]string,256)
		i := 0
		for {
			msg := <-b.Messages
			content,left := parseMultiPart(msg)
			parts[i] = content
			i++
			if left == 0 {
				break
			}
		}
		var content string
		content = strings.Join(parts,"")

		msg,err := tvpn.ParseMessage(content)
		if err == nil {
			return msg
		}
	}
}


