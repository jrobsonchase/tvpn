package connect

import (
	"net"
	"log"
	//"bufio"
)


func ServerConnect(host, port string) (convo *CConv, err error) {

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	convo = initConvo(host,conn)

	return convo, nil
}

func ConverseClient(convo *CConv,name string) {
	convo.To <- ConnMsg{Type: Hello, To: convo.With, From: name}

	msg := <-convo.From
	if msg.Type != Hello {
		log.Print("Server failed to send HELLO!")
		convo.End()
		convo.conn.Close()
		return
	}
	convo.To <- ConnMsg{Type: Bye, To: convo.With, From: name}
	convo.End()
	return
}
