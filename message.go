package tvpn

import (
	"fmt"
	"regexp"
)

type Message struct {
	To, From string
	Type     int
	Data     map[string]string
}

const (
	Init int = iota
	Join
	Quit
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
	initRE     string = `^INIT$`
	acceptRE          = `^ACCEPT$`
	denyRE            = `^DENY (?P<reason>.*)$`
	dhpubRE           = `^DHPUB (?P<x>[A-Za-z0-9+/=]+) (?P<y>[A-Za-z0-9+/=]+)$`
	tunnipRE          = `^TUNNIP (?P<ip>[0-9]{1,3}(?:\.[0-9]{1,3}){3})$`
	conninfoRE        = `^CONNINFO (?P<ip>[0-9]{1,3}(?:\.[0-9]{1,3}){3}) (?P<port>[0-9]+)$`
)

func (m Message) Print() {
	fmt.Printf("From: %s\nTo: %s\nType: %d\n", m.From, m.To, m.Type)
	for i, v := range m.Data {
		fmt.Printf("%s: %s\n", i, v)
	}
}

func ParseMessage(message string) (*Message, error) {

	init := regexp.MustCompile(initRE)
	accept := regexp.MustCompile(acceptRE)
	deny := regexp.MustCompile(denyRE)
	dhpub := regexp.MustCompile(dhpubRE)
	tunnip := regexp.MustCompile(tunnipRE)
	conninfo := regexp.MustCompile(conninfoRE)

	var data map[string]string = make(map[string]string)

	switch {
	case init.MatchString(message):
		return &Message{Type: Init, Data: data}, nil

	case accept.MatchString(message):
		return &Message{Type: Accept, Data: data}, nil

	case deny.MatchString(message):
		data["reason"] = deny.ReplaceAllString(message, "${reason}")
		return &Message{Type: Deny, Data: data}, nil

	case dhpub.MatchString(message):
		data["x"] = dhpub.ReplaceAllString(message, "${x}")
		data["y"] = dhpub.ReplaceAllString(message, "${y}")
		return &Message{Type: Dhpub, Data: data}, nil

	case tunnip.MatchString(message):
		data["ip"] = tunnip.ReplaceAllString(message, "${ip}")
		return &Message{Type: Tunnip, Data: data}, nil

	case conninfo.MatchString(message):
		data["ip"] = conninfo.ReplaceAllString(message, "${ip}")
		data["port"] = conninfo.ReplaceAllString(message, "${port}")
		return &Message{Type: Conninfo, Data: data}, nil

	default:
		return nil, messageError{fmt.Sprintf("Failed to parse message: %s", message)}
	}

}

func (m Message) String() string {
	switch m.Type {
	case Init:
		return "INIT"
	case Deny:
		return fmt.Sprintf("DENY %s", m.Data["reason"])
	case Accept:
		return "ACCEPT"
	case Dhpub:
		return fmt.Sprintf("DHPUB %s %s", m.Data["x"], m.Data["y"])
	case Tunnip:
		return fmt.Sprintf("TUNNIP %s", m.Data["ip"])
	case Conninfo:
		return fmt.Sprintf("CONNINFO %s %s", m.Data["ip"], m.Data["port"])
	}
	return ""
}
