package main

import (
	"tvpn"
	"tvpn/irc"
	"fmt"
	//"tvpn/ovpn"
	//"tvpn/stun"
	"os"
)

func main() {
	irc, err := irc.Connect("chat.freenode.net:6667","hurdurclient","#testchan412")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	for {
		msg := irc.RecvMessage()
		if msg.Type == tvpn.Join {
			irc.SendMessage(tvpn.Message{Type: tvpn.Init, To: msg.From})
			msg := irc.RecvMessage()
			if msg.Type == tvpn.Accept {
				fmt.Printf("Accepted! Yay!\n")
			}
		}
	}
}


