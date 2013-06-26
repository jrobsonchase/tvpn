package connect

import (
	"strings"
	"bufio"
	"log"
	"io"
	"fmt"
	"net"
	"runtime"
)

type ConnMsg struct {
	Type string
	IP net.IP
	Port int
	Key string
}

type CConv struct {
	Msgs chan ConnMsg
	Repl chan ConnMsg
	done chan bool
}

const (
	Init string = "INIT"
	AccInfo string = "ACCINFO"
	Acc string = "ACC"
	Rej string = "REJ"
)

func parseString(conn io.Reader,out *string) error {
	_,err := fmt.Fscanf(conn,"%s",out)
	if err != nil {
		log.Printf("Failed to parse string!")
		return err
	}
	return nil
}

func parseInt(conn io.Reader,out *int) error {
	_,err := fmt.Fscanf(conn,"%d",out)
	if err != nil {
		log.Printf("Failed to parse int!")
		return err
	}
	return nil
}

func parseInfo(conn io.Reader,sender,ip *string,port *int,key *string) error {
	err := parseString(conn,sender)
	if err != nil {
		return err
	}
	err = parseString(conn,ip)
	if err != nil {
		return err
	}
	err = parseInt(conn,port)
	if err != nil {
		return err
	}
	err = parseString(conn,key)
	if err != nil {
		return err
	}
	return nil
}

func sendMessage(conn io.Writer, recip string, msg ConnMsg) error {
	switch msg.Type {
	case Init:
		fmt.Fprintf(conn,"%s %s %s %d %s",recip,Init,msg.IP.String(),msg.Port,msg.Key)
	case AccInfo:
		fmt.Fprintf(conn,"%s %s %s %d %s",recip,AccInfo,msg.IP.String(),msg.Port,msg.Key)
	case Acc:
		fmt.Fprintf(conn,"%s %s",recip,Acc)
	case Rej:
		fmt.Fprintf(conn,"%s %s",recip,msg.Key)
	default:
		log.Print("Invalid message type!")
	}
	return nil
}

func DispatchMsgs(conn bufio.ReadWriter) error {
	convos := make(map[string] CConv)

	// Oh god this is ugly.
	// creates a routine that continually queries the current
	// conversations for messages or disconnects
	// probably broken
	go func(cs *map[string] CConv) {
		for {
			l := len(*cs)
			s := make(chan bool,l)
			for i,v := range *cs {
				go func(n string,c CConv,sem chan bool) {
					select {
					case msg := <-c.Repl:
						sendMessage(conn,i,msg)
					case <-c.done:
						delete(*cs,n)
					}
					sem <- true
				}(i,v,s)
			}
			for i := 0; i < l; i++ {
				<-s
			}
			runtime.Gosched()
		}
	}(&convos)


	log.Print("begin message handling...")

	for i := 0; i < 128; {
		line,err := conn.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read line: %s",err)
		}
		inStrs := make([]string, 10)
		scanner := bufio.NewScanner(strings.NewReader(line))
		scanner.Split(bufio.ScanWords)

		for j := 0; j < 10 && scanner.Scan(); {
			inStrs[j] = scanner.Text()
			log.Printf("Read %s from net",inStrs[j])
			j++
		}

		// This code is bad
		// And I should feel bad
		switch inStrs[0] {
			case Init:
				name := inStrs[1]
				ip := inStrs[2]
				var port int
				fmt.Fscanf(strings.NewReader(inStrs[3]),"%d",&port)
				key := inStrs[4]
				msg := ConnMsg{Type: Init, IP: net.ParseIP(ip),Port: port, Key: key}
				convos[name] = CConv{Msgs: make(chan ConnMsg),Repl: make(chan ConnMsg),done: make(chan bool)}
				go converse(convos[name])
				convos[name].Msgs <- msg
				i = len(convos)
			case AccInfo:
				name := inStrs[1]
				// need break here if convo doesn't exist
				ip := inStrs[2]
				var port int
				fmt.Fscanf(strings.NewReader(inStrs[3]),"%d",&port)
				key := inStrs[4]
				msg := ConnMsg{Type: AccInfo, IP: net.ParseIP(ip),Port: port, Key: key}
				convos[name].Msgs <- msg
			case Acc:
				name := inStrs[1]
				// and here
				msg := ConnMsg{Type: Acc}
				convos[name].Msgs <- msg
			case Rej:
				name := inStrs[1]
				// and here
				reason := inStrs[2]
				msg := ConnMsg{Type: Rej, Key: reason}
				convos[name].Msgs <- msg
			default:
				log.Printf("TypeString not recognized: %s",inStrs[0])
		}
	}
	return nil
}

func converse(conv CConv) {
	log.Print("Starting new convo...")
	defer func() {conv.done <- true}()
	for {
		msg := <-conv.Msgs
		log.Printf("%s %s %d %s",msg.Type,msg.IP,msg.Port,msg.Key)
		conv.Repl <-ConnMsg{Type: AccInfo,IP: net.ParseIP("0.0.0.0"),Port: 25, Key: "1234"}
	}
}

