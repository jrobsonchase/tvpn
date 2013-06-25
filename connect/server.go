package connect

import (
	"log"
	"bytes"
	"io"
	"fmt"
	"tvpn/stun"
)

func ListenReq(name string,conn io.ReadWriter) error {
	buff := new(bytes.Buffer)
	for {
		_,err := buff.ReadFrom(conn)
		if err != nil {
			return err
		}
		go handleInput(buff,conn,name)
	}
}

func handleInput(buff *bytes.Buffer,conn io.ReadWriter,name string) {
	req := new(ConnReq)

	_,err := fmt.Fscanf(buff,"%s %s %d %s",&req.Name,&req.IP,&req.Port,&req.PubKey)
	if err != nil {
		log.Printf("Failed to parse connection request: %s",err)
		return
	}

	//err := pki.ValidateKey(req.PubKey)

	external := stun.DiscoverExternal(int(rgen.Uint32() % (65535 - 1024) + 1024),StunServer)
	if external == nil {
		log.Print("Failed to get external address")
		return
	}

	resp,err := SendRequest(conn,connReqResp,req.Name,name,external.IP,external.Port,PubKey)
	if err != nil {
		log.Printf("Failed to complete handshake: %s",err)
	}
	if resp[:7] == connRefuse[:7] {
		log.Printf("Remote host refused connection: %s",resp[8:])
		return
	}
	if resp == connAccept {
		/* Connection happens here */
	}
}

