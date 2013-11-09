package main

import (
	"fmt"
	"os"
	"github.com/Pursuit92/tvpn"
	"github.com/Pursuit92/tvpn/irc"
)

func main() {
	irc.SetLogLevel(5)
	i, err := irc.Connect("chat.freenode.net:6667", "hurdurtestnick", "#joshtestgroup")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	for {
		msg := i.RecvMessage()
		fmt.Printf("Received: %s\n",msg.String())
		if msg.Type == tvpn.Init {
			fmt.Println("Got INIT")
			data := make(map[string]string)
			data["reason"] = "You're ugly"
			i.SendMessage(tvpn.Message{Type: tvpn.Deny, To: msg.From, Data: data})
		} else if msg.Type == tvpn.Join {
			fmt.Println("Got Join")
			i.SendMessage(tvpn.Message{Type: tvpn.Init, To: msg.From})
		}
	}
}
