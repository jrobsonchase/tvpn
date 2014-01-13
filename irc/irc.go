/*
 *  TVPN: A Peer-to-Peer VPN solution for traversing NAT firewalls
 *  Copyright (C) 2013  Joshua Chase <jcjoshuachase@gmail.com>
 *
 *  This program is free software; you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation; either version 2 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License along
 *  with this program; if not, write to the Free Software Foundation, Inc.,
 *  51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/

package irc

import (
	"sync"
	"github.com/Pursuit92/irc"
	"github.com/Pursuit92/tvpn"
)

func SetLogLevel(n int) {
	irc.SetLogLevel(n)
}

type IRCBackend struct {
	Nick,Chan,Server string
	Conn        *irc.Conn
	Messages    chan irc.CmdErr
}

func (i *IRCBackend) Configure(conf tvpn.SigConfig) {
	i.Nick = conf["Name"]
	i.Chan = conf["Group"]
	i.Server = conf["Server"]

	if i.Conn != nil {
		// cleanup old connection stuff
	}

	err := i.Connect()

	if err != nil {
		panic(err)
	}

}

func (i *IRCBackend) Connect() error {

	conn, err := irc.DialIRC(i.Server, []string{i.Nick}, i.Nick, i.Nick)
	if err != nil {
		return err
	}
	_, err = conn.Register()
	if err != nil {
		return err
	}

	chann, err := conn.Join(i.Chan)
	if err != nil {
		return err
	}

	joinpart, _ := chann.Expect(irc.Command{"", "(JOIN)|(PART)", []string{}})
	quit, _ := conn.Expect(irc.Command{"", "QUIT", []string{}})

	msgs, err := conn.Expect(irc.Command{"", "PRIVMSG", []string{i.Nick,".*"}})
	if err != nil {
		return err
	}

	i.Messages = make(chan irc.CmdErr)

	combine([]<-chan irc.CmdErr{joinpart.Chan,quit.Chan,msgs.Chan},i.Messages)

	//users := chann.GetUsers()

	i.Conn = conn

	return nil
}


func (b IRCBackend) RecvMessage() (tvpn.Message,error) {
	for {
		input := <-b.Messages
		switch input.Cmd.Command {
		case "QUIT", "PART":
			return tvpn.Message{From: input.Cmd.Message().Nick, Type: tvpn.Quit}, nil
		case "JOIN":
			if input.Cmd.Message().Nick != b.Conn.Nick {
				return tvpn.Message{From: input.Cmd.Message().Nick, Type: tvpn.Join}, nil
			}
		default:
			ircMsg := input.Cmd.Message()
			msg, err := tvpn.ParseMessage(input.Cmd.Params[len(input.Cmd.Params)-1])
			if err == nil {
				msg.From = ircMsg.Nick
				return *msg, nil
			}
		}
	}
}

func (b IRCBackend) SendMessage(mes tvpn.Message) error {
	return b.Conn.Send(irc.Command{b.Conn.Nick, irc.Privmsg, []string{mes.To, mes.String()}})
}

func combine(inputs []<-chan irc.CmdErr, output chan<- irc.CmdErr) {
	var group sync.WaitGroup
	for i := range inputs {
		group.Add(1)
		go func(input <-chan irc.CmdErr) {
			for val := range input {
				output <- val
			}
			group.Done()
		} (inputs[i])
	}
	go func() {
		group.Wait()
		close(output)
	} ()
}
