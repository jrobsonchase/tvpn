package main

import (
	"bufio"
	"flag"
	"os"
	"fmt"
	"strings"
	"tvpn/tvpn"
)

const friendLimit int = 256

func exitError(s string) {
	fmt.Printf("%s\n",s)
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
	flag.Parse()

	var friends []string

	if *friendFile != "" {
		file,err := os.Open(*friendFile)
		if err != nil {
			fmt.Printf("Error openfing file: %s\n",err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(file)
		friends = make([]string,friendLimit)
		var i int
		for scanner.Scan() {
			friends[i] = strings.Trim(scanner.Text(),"\t \n")
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

	fmt.Printf("Loaded friends:\n")
	for _,v := range friends {
		fmt.Println(v)
	}

	tvpnInstance := tvpn.TVPN{
		Name: *ircNick,
		Group: *ircChannel,
		IRCServer: *ircString,
		STUNServer: *stunString,
	}

	err := tvpnInstance.Run()

}
