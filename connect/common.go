package connect

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	//	"runtime"
)

type ConnMsg struct {
	Type string
	From string
	To   string
	IP   net.IP
	Port int
	Key  string
}

type MsgErr struct {
	Text string
}

func (m *MsgErr) Error() string {
	return "Error parsing token: " + m.Text
}

type CConv struct {
	With string
	To   chan ConnMsg
	From chan ConnMsg
	Done chan bool
	mut chan bool
	conn net.Conn
}

type ConnErr string

func (c *ConnErr) Error() string {
	msg := string(*c) + "did not send HELLO!"
	return msg
}

const (
	Hello   string = "HELLO"
	Bye		string = "BYE"
	Init    string = "INIT"
	AccInfo string = "ACCINFO"
	Acc     string = "ACC"
	Rej     string = "REJ"
)

func validType(s string) bool {
	switch s {
	case Hello,Bye,Init,AccInfo,Acc,Rej:
		return true
	default:
		return false
	}
}

func readMsgs(from *bufio.Reader, out chan ConnMsg, done chan bool) {
	for {
		msg, err := readMessage(from)
		if err != nil {
			break
		}
		out <- *msg
	}
	done <-true
}

func readMessage(from *bufio.Reader) (msg *ConnMsg, err error) {
	line, err := from.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read line: %s", err)
		return nil, err
	}

	input := make([]string, 10)
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)

	for j := 0; j < 10 && scanner.Scan(); {
		input[j] = scanner.Text()
		log.Printf("Read %s from net", input[j])
		j++
	}

	to := input[0]
	inStrs := input[1:]

	var name string
	// make sure the command is valid, otherwise restart loop
	switch {
	case validType(inStrs[0]):
		name = inStrs[1]
		msg = &ConnMsg{Type: inStrs[0],To: to, From: name}
	default:
		return nil, &MsgErr{inStrs[0]}
	}

	// Grab other info from the messages
	switch inStrs[0] {
	case Init, AccInfo:
		ip := inStrs[2]
		var port int
		fmt.Fscanf(strings.NewReader(inStrs[3]), "%d", &port)
		key := inStrs[4]
		msg.IP = net.ParseIP(ip)
		msg.Port = port
		msg.Key = key
	case Rej:
		reason := inStrs[2]
		msg.Key = reason
	}

	// finally, send the message to the convo
	log.Printf("returning %s from %s",msg.Type,msg.From)
	return msg, nil
}

func sendMessage(conn *bufio.Writer, msg ConnMsg) error {
	var err error
	if validType(msg.Type) {
		_,err = fmt.Fprintf(conn,"%s %s %s ",msg.To,msg.Type,msg.From)
	} else {
		log.Print("Invalid message type!")
		return &MsgErr{msg.Type}
	}

	switch msg.Type {
	case Init,AccInfo:
		_,err = fmt.Fprintf(conn, "%s %d %s\n", msg.IP.String(), msg.Port, msg.Key)
	case Rej:
		_,err = fmt.Fprintf(conn, "%s\n", msg.To, Rej, msg.Key)
	default:
		_,err = fmt.Fprintf(conn,"\n")
	}
	if err != nil {
		return err
	}
	conn.Flush()
	log.Print("Sent Message!")
	return nil
}

func sendMsgs(write *bufio.Writer, msgs chan ConnMsg,done chan bool) {
	for {
		msg,ok := <-msgs
		if ! ok {
			break
		}
		err := sendMessage(write, msg)
		if err != nil {
			log.Print(err)
			done <-true
			break
		}
	}
}

func initConvo(name string,conn net.Conn) *CConv {
	read := bufio.NewReader(conn)
	write := bufio.NewWriter(conn)
	toChan := make(chan ConnMsg)
	fromChan := make(chan ConnMsg)
	done := make(chan bool,2)
	go readMsgs(read, fromChan, done)
	go sendMsgs(write, toChan,done)

	return &CConv{With: name, From: fromChan, To: toChan, Done: done,conn: conn}
}

func (convo *CConv) End() {
	close(	convo.To)
	convo.conn.Close()
}

func printMsg(msg ConnMsg) {
	fmt.Printf("%s %s %s %s %d %s",msg.To,msg.Type,msg.From,msg.IP,msg.Port,msg.Key)
}
