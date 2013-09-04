package tvpn

import (
	"regexp"
	"fmt"
)

type Message struct {
	To,From string
	Type int
	Data map[string] string
}

const (
	Init int = iota
	Join
	Accept
	Deny
	Dhpub
	Tunnip
	Conninfo
)

type messageError struct {
	message string
}

func (e messageError) Error() string {
	return e.message
}

const (
	initRE string = `^INIT$`
	acceptRE = `^ACCEPT$`
	denyRE = `^DENY (?P<reason>.*)$`
	dhpubRE = `^DHPUB (?P<x>[a-f0-9]+) (?P<y>[a-f0-9]+)$`
	tunnipRE = `^TUNNIP (?P<ip>[0-9]{1,3}(?:\.[0-9]{1,3}){3})$`
	conninfoRE = `^CONNINFO (?P<ip>[0-9]{1,3}(?:\.[0-9]{1,3}){3}) (?P<port>[0-9]+)$`
)

func ParseMessage(message string) (*Message,error) {

	init := regexp.MustCompile(initRE)
	accept := regexp.MustCompile(acceptRE)
	deny := regexp.MustCompile(denyRE)
	dhpub := regexp.MustCompile(dhpubRE)
	tunnip := regexp.MustCompile(tunnipRE)
	conninfo := regexp.MustCompile(conninfoRE)

	var data map[string] string = make(map[string] string)

	switch {
	case init.MatchString(message):
		return &Message{Type: Init,Data: data},nil

	case accept.MatchString(message):
		return &Message{Type: Accept,Data: data},nil

	case deny.MatchString(message):
		data["reason"] = deny.ReplaceAllString(message,"${reason}")
		return &Message{Type: Deny,Data: data},nil

	case dhpub.MatchString(message):
		data["x"] = dhpub.ReplaceAllString(message,"${x}")
		data["y"] = dhpub.ReplaceAllString(message,"${y}")
		return &Message{Type: Dhpub,Data: data},nil

	case tunnip.MatchString(message):
		data["ip"] = conninfo.ReplaceAllString(message,"${ip}")
		return &Message{Type: Tunnip,Data: data},nil

	case conninfo.MatchString(message):
		data["ip"] = conninfo.ReplaceAllString(message,"${ip}")
		data["port"] = conninfo.ReplaceAllString(message,"${port}")
		return &Message{Type: Conninfo,Data: data},nil

	default:
		return nil,messageError{fmt.Sprintf("Failed to parse message: %s",message)}
	}

}
