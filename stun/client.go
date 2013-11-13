package stun

import (
	"encoding/binary"
	"net"
	"github.com/Pursuit92/LeveledLogger/log"
)

type StunErr string

func (s StunErr) Error() string {
	return string(s)
}

type StunBackend struct {
	Server string
}

func (b StunBackend) DiscoverExt(port int) (net.IP,int,error) {
	uadd := DiscoverExternal(port,b.Server)
	if uadd == nil {
		return nil,0,StunErr("Failed to get external address!")
	}
	return uadd.IP,uadd.Port,nil
}

func DiscoverExternal(port int, addr string) (*net.UDPAddr) {

	// Make the initial connection to the helper server
	rem := "udp"
	laddr := &net.UDPAddr{Port: port}
	helper,err := net.ResolveUDPAddr(rem,addr)
	if err != nil {
		log.Out.Printf(2,"Failed to resolve UDP: %s",err)
		return nil
	}
	conn, err := net.DialUDP(rem, laddr, helper)
	if err != nil {
		log.Out.Printf(2,"Failed to bind UDP: %s",err)
		return nil
	}
	defer conn.Close()

	// send the request for info
	n, err := conn.Write([]byte{Req})
	if err != nil {
		log.Out.Printf(2,"Failed to send request: %s",err)
		return nil
	}
	log.Out.Printf(3,"Wrote %d bytes to %s",n,rem)

	// server response - 1 byte for type, 4 for the port, and 16 for the IP
	resp := make([]byte, 25)
	n, err = conn.Read(resp)
	if err != nil {
		log.Out.Printf(2,"Failed to read response: %s",err)
		return nil
	}
	log.Out.Printf(3,"Read %d bytes from %s",n,rem)

	// slice out the port bytes
	portBytes := resp[1:9]

	// test for ipv6 based on number of bytes sent and slice the ip bytes
	var ipBytes []byte
	if n > 10 {
		ipBytes = resp[9:]
	} else {
		ipBytes = resp[9:13]
	}

	// convert the port bytes to int64
	extPort,n := binary.Varint(portBytes)
	if n == 0 {
		log.Out.Print(2,"Failed to read integer from portBytes")
		return nil
	}


	return &net.UDPAddr{ IP: ipBytes, Port: int(extPort) }

}

func SetLogLevel(n int) {
	log.Out.SetLevel(n)
}

func SetLogPrefix(s string) {
	log.Out.SetPrefix(2,s)
	log.Out.SetPrefix(3,s)
}
