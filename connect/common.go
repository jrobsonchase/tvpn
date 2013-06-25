package connect

import (
	"io"
	"fmt"
	"bytes"
	"net"
)

type ConnMsg struct {
	Type byte
	IP []byte
	Port int
	PubKey string
}

type CConv struct {
	Name string
	Pipe [2]chan ConnMsg
}

const (
	Req byte = iota
	ReqResp
	Accept
	Refuse
)

func parseInfo(conn io.Reader,*sender,*ip string,port *int,key *string) error {
	_,err := fmt.Fscanf(conn,"%s %s %d %s",sender,ip,port,key)
	if err != nil {
		log.Printf("Failed to parse message info!")
		return err
	}
}

func DispatchMsgs(conn io.ReadWriter) error {
	convos := make([]CConv, 128)
	var typeStr string
	var name string
	var ip string
	var port int
	var key string

	for i := 0; i < 128; {
		_,err := fmt.Fscanf(conn,"%s", &typeStr)
		if err != nil {
			log.Printf("Failed to parse message command!")
			continue
		}

		if typeStr == InitStr {
			parseInfo(conn,&name,&ip,&port,&key)
			msg := ConnMsg{Type: Init, IP: net.ParseIP(ip),Port: port, PubKey: key}
			// TODO: This is where I stopped
		}
	}
}

