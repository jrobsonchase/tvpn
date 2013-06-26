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

	// creates a routine that continually queries the current
	// conversations for messages or disconnects
	go manageConvos(conn, &convos)


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

		var name string
		var msg ConnMsg
		// make sure the command is valid, otherwise restart loop
		switch inStrs[0] {
		case Init,AccInfo,Acc,Rej:
			name = inStrs[1]
		default:
			log.Printf("TypeString not recognized: %s",inStrs[0])
			continue
		}

		// if it's a message for an existing convo, make sure it's been Init'd
		switch inStrs[0] {
		case AccInfo,Acc,Rej:
			_,exists := convos[name]
			if ! exists {
				log.Printf("Conversation with %s does not exist!",name)
				continue
			}
		default:
		}

		// Grab other info from the messages
		switch inStrs[0] {
		case Init,AccInfo:
			ip := inStrs[2]
			var port int
			fmt.Fscanf(strings.NewReader(inStrs[3]),"%d",&port)
			key := inStrs[4]
			msg = ConnMsg{Type: inStrs[0], IP: net.ParseIP(ip),Port: port, Key: key}
		case Acc:
			msg = ConnMsg{Type: Acc}
		case Rej:
			reason := inStrs[2]
			msg = ConnMsg{Type: Rej, Key: reason}
		}

		// if it's init, a new convo needs to start
		if inStrs[0] == Init {
			convos[name] = CConv{Msgs: make(chan ConnMsg),Repl: make(chan ConnMsg),done: make(chan bool)}
			go converse(name,convos[name])
			i = len(convos)
		}

		// finally, send the message to the convo
		convos[name].Msgs <- msg
	}
	return nil
}

// this is a stub. need to put real logic and functionality into it
func converse(name string,conv CConv) {
	log.Print("Starting new convo...")
	defer func() {conv.done <- true}()
	for {
		msg := <-conv.Msgs
		log.Printf("%s %s %d %s",msg.Type,msg.IP,msg.Port,msg.Key)
		conv.Repl <-ConnMsg{Type: AccInfo,IP: net.ParseIP("0.0.0.0"),Port: 25, Key: "1234"}
	}
}

// repeatedly query the conversations for messages/done notifications
func manageConvos(conn bufio.ReadWriter,cs *map[string] CConv) {
	for {
		l := len(*cs)
		// semaphore to make sure we only restart the iteration
		// when all convos have had a chance to go
		s := make(chan bool,l)
		// fire off all of the convos
		for i,v := range *cs {
			go func(n string,c CConv,sem chan bool) {
				select {
				case msg := <-c.Repl:
					sendMessage(conn,i,msg)
				case <-c.done:
					delete(*cs,n)
				default:
				}
				// say that it's done with this go-round
				sem <- true
			}(i,v,s)
		}
		// wait for everyone to finish
		for i := 0; i < l; i++ {
			<-s
		}
		// let someone else go - may not be necessary
		runtime.Gosched()
	}
}
