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
	"github.com/Pursuit92/LeveledLogger/log"
)

func SetLogLevel(n int) {
	irc.SetLogLevel(n)
	log.Out.SetLevel(n)
}

type IRCBackend struct {
	Nick,Chan,Server string
	Conn        *irc.Conn
	Messages    chan irc.CmdErr
}

func (i *IRCBackend) Configure(conf tvpn.SigConfig) {
	var changed, exists bool
	if i.Nick != conf["Name"] || i.Chan != conf["Group"] || i.Server != conf["Server"] {
		i.Nick = conf["Name"]
		i.Chan = conf["Group"]
		i.Server = conf["Server"]
		changed = true
	}

	if i.Conn != nil {
		exists = true
	}

	if changed {
		if exists {
			i.Disconnect()
		}
		err := i.Connect()

		if err != nil {
			log.Out.Println(1,err)
		}
	}
}

func (i *IRCBackend) Disconnect() {
	i.Conn.Quit()
}

func (i *IRCBackend) Reconnect() error {
	i.Disconnect()
	return i.Connect()
}

func (i *IRCBackend) Connect() error {

	log.Out.Printf(2,"Connecting to irc server %s with nick %s\n",i.Server, i.Nick)

	conn, err := irc.DialIRC(i.Server, []string{i.Nick}, i.Nick, i.Nick)
	if err != nil {
		return err
	}

	log.Out.Println(2,"Connect success, attempting register...")

	_, err = conn.Register()
	if err != nil {
		return err
	}

	log.Out.Printf(2,"Register success, attempting to join channel %s\n",i.Chan)

	chann, err := conn.Join(i.Chan)
	if err != nil {
		return err
	}

	log.Out.Println(2,"Join success, registering listeners for join/part")
	joinpart, _ := chann.Expect(irc.Command{"", "(JOIN)|(PART)", []string{}})
	quit, _ := conn.Expect(irc.Command{"", "QUIT", []string{}})

	msgs, err := conn.Expect(irc.Command{"", "PRIVMSG", []string{i.Nick,".*"}})
	if err != nil {
		return err
	}

	i.Messages = make(chan irc.CmdErr)

	combine([]<-chan irc.CmdErr{joinpart.Chan,quit.Chan,msgs.Chan},i.Messages)

	i.Conn = conn

	log.Out.Println(2,"Connect complete!")
	return nil
}


func (b IRCBackend) RecvMessage() (tvpn.Message,error) {
	for input := range b.Messages {
		if input.Err != nil {
			return tvpn.Message{},input.Err
		}
		switch input.Cmd.Command {
		case "QUIT", "PART":
			log.Out.Printf(2,"Received QUIT/PART from %s\n",input.Cmd.Message().Nick)
			return tvpn.Message{From: input.Cmd.Message().Nick, Type: tvpn.Quit}, nil
		case "JOIN":
			if input.Cmd.Message().Nick != b.Conn.Nick {
				log.Out.Printf(2,"Received JOIN from %s\n",input.Cmd.Message().Nick)
				return tvpn.Message{From: input.Cmd.Message().Nick, Type: tvpn.Join}, nil
			}
		default:
			ircMsg := input.Cmd.Message()
			msg, err := tvpn.ParseMessage(input.Cmd.Params[len(input.Cmd.Params)-1])
			if err == nil {
				log.Out.Printf(2,"Received message: %s\n",input.Cmd.String())
				msg.From = ircMsg.Nick
				return *msg, nil
			} else {
				log.Out.Printf(2,"Received malformed message: %s\n",input.Cmd.String())
			}
		}
	}
	return tvpn.Message{}, irc.IRCErr("Disconnected")
}

func (b IRCBackend) SendMessage(mes tvpn.Message) error {
	log.Out.Printf(2,"Sending message to %s: %s\n",mes.To,mes.String())
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
