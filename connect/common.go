package connect

import (
	"io"
	"fmt"
)
type ConnReq struct {
	Name string
	IP []byte
	Port int
	PubKey string
}

const (
	reqInfo string = "%s %s %d %s"
	connReq string = "%s CONNREQ " + reqInfo
	connReqResp string = "%s CONNREQRESP " + reqInfo
	connAccept string = "CONNACC"
	connRefuse string = "CONNREF %s"
)

func SendRequest(conn io.ReadWriter,com string, info ...interface{}) (string,error) {
	fmt.Fprintf(conn,com,info...)
	return "",nil
}
