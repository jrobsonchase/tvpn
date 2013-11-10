package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/Pursuit92/tvpn"
	"github.com/Pursuit92/tvpn/ovpn"
	"github.com/Pursuit92/tvpn/stun"
	"github.com/Pursuit92/tvpn/irc"
)

const friendLimit int = 256

func exitError(s string) {
	fmt.Printf("%s\n", s)
	os.Exit(1)
}

func main() {
	friendFile := flag.String("friends", "", "File containing friends. One per line")
	ircChannel := flag.String("group", "", "Connection group to join. Ensures that you get updates on your friends' presence")
	ircString := flag.String("irc", "irc.freenode.net:6667", "IRC server info")
	//ircConn := flag.Bool("tls", false, "Use TLS for IRC connection? (unimplemented)")
	ircNick := flag.String("name", "", "Name to use when connecting to the IRC server")
	//ircPass := flag.String("pass", "", "Optional password for IRC connection")
	//ircIdent := flag.String("identify", "", "Optional password for NickServ identification")
	stunString := flag.String("stun", "", "STUN server info")
	debugLevel := flag.Int("d",1,"Debugging level. Set to 1 by default")
	flag.Parse()

	var friends []string

	if *friendFile != "" {
		file, err := os.Open(*friendFile)
		if err != nil {
			fmt.Printf("Error openfing file: %s\n", err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(file)
		friends = make([]string, friendLimit)
		var i int
		for scanner.Scan() {
			friends[i] = strings.Trim(scanner.Text(), "\t \n")
			i++
		}
		friends = friends[:i]

	} else {
		exitError("You must specify a friends file with -friends")
	}

	if *ircChannel == "" {
		exitError("You must specify a group to join with -group")
	}

	if *ircNick == "" {
		exitError("You must specify an IRC name with -name")
	}

	if *stunString == "" {
		exitError("You must specify a STUN server with -stun")
	}

	switch *debugLevel {
	case 0,1:
	case 2:
		irc.SetLogLevel(2)
		tvpn.SetLogLevel(2)
	case 3:
		stun.SetLogLevel(2)
		irc.SetLogLevel(2)
		tvpn.SetLogLevel(2)
	case 4:
		stun.SetLogLevel(2)
		irc.SetLogLevel(2)
		tvpn.SetLogLevel(3)
	case 5:
		stun.SetLogLevel(3)
		irc.SetLogLevel(3)
		tvpn.SetLogLevel(3)
	case 6,7,8,9,10:
		stun.SetLogLevel(10)
		irc.SetLogLevel(10)
		tvpn.SetLogLevel(10)
	default:
	}

	fmt.Printf("Loaded friends:\n")
	for _, v := range friends {
		fmt.Println(v)
	}

	ircBackend, err := irc.Connect(*ircString, *ircNick, *ircChannel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tvpnInstance := tvpn.New(
		*ircNick,
		*ircChannel,
		friends,
		ircBackend,
		stun.StunBackend{*stunString},
		ovpn.New(),
		tvpn.NewIPManager("3.0.0.0",256))

	err = tvpnInstance.Run()

}
