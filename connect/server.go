package connect

import (
	"bufio"
	"log"
	"net"
)

type ConnServer struct {
	net.Listener
	Name string
}

func StartServer(name,port string,convos chan CConv) error {
	ln,err := net.Listen("tcp",":"+port)
	if err != nil {
		return err
	}

	server := ConnServer{ln,name}

	for {
		conv,err := server.AcceptClient()
		if err != nil {
			log.Printf("Server error: %s",err)
		}
		convos <- *conv
	}
}


func (serv ConnServer) AcceptClient() (*CConv, error) {
	conn, err := serv.Accept()
	if err != nil {
		log.Printf("Accept Error: %s", err)
		return nil, err
	}

	connRead := bufio.NewReader(conn)

	msg, err := readMessage(connRead)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if msg.Type != Hello {
		e := ConnErr(msg.From)
		conn.Close()
		return nil, &e
	}

	convo := initConvo(msg.From,conn)

	log.Printf("retrning convo")
	return convo, nil
}

func ConverseServer(convo *CConv,name string) {
	convo.To <- ConnMsg{Type: Hello, To: convo.With, From: name}
	for {
		msg := <-convo.From
		if msg.Type == Bye {
			log.Print("Client disconnected!")
			convo.End()
			break
		}
		printMsg(msg)
	}
}
