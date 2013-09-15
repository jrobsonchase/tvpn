package stun

import (
	"encoding/binary"
	"net"
	"log"
)

type StunErr string

func (s StunErr) Error() string {
	return string(s)
}

type StunBackend struct {
	server string
}

func (b StunBackend) DiscoverExt(port int) (net.IP,int,error) {
	uadd := DiscoverExternal(port,b.server)
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
		log.Printf("Failed to resolve UDP: %s",err)
		return nil
	}
	conn, err := net.DialUDP(rem, laddr, helper)
	if err != nil {
		log.Printf("Failed to bind UDP: %s",err)
		return nil
	}

	// send the request for info
	n, err := conn.Write([]byte{Req})
	if err != nil {
		log.Printf("Failed to send request: %s",err)
		return nil
	}
	log.Printf("Wrote %d bytes to %s",n,rem)

	// server response - 1 byte for type, 4 for the port, and 16 for the IP
	resp := make([]byte, 25)
	n, err = conn.Read(resp)
	if err != nil {
		log.Printf("Failed to read response: %s",err)
		return nil
	}
	log.Printf("Read %d bytes from %s",n,rem)

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
		log.Print("Failed to read integer from portBytes")
		return nil
	}

	return &net.UDPAddr{ IP: ipBytes, Port: int(extPort) }

}




