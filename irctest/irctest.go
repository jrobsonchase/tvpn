package main

import (
	"fmt"
	"os"
	"github.com/Pursuit92/tvpn"
	"github.com/Pursuit92/tvpn/irc"
)

func main() {
	irc, err := irc.Connect("chat.freenode.net:6667", "hurdurtestnick", "#testchan412")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	for {
		msg := irc.RecvMessage()
		if msg.Type == tvpn.Init {
			data := make(map[string]string)
			data["reason"] = "You're ugly"
			irc.SendMessage(tvpn.Message{Type: tvpn.Deny, To: msg.From, Data: data})
		}
	}
}
