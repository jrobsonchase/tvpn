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
	irc, err := irc.Connect("chat.freenode.net:6667","hurdurserver","#testchan412")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	for {
		msg := irc.RecvMessage()
		if msg.Type == tvpn.Init {
			irc.SendMessage(tvpn.Message{Type: tvpn.Accept, To: msg.From})
		}
	}
}
