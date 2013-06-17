package stun

import (
	"encoding/binary"
	"net"
	"log"
)

func DiscoverExternal(port int,helper *net.UDPAddr) (*net.UDPAddr) {

	// Make the initial connection to the helper server
	rem := "udp"
	laddr := &net.UDPAddr{Port: port}
	conn, err := net.DialUDP(rem, laddr, helper)
	if err != nil {
		log.Fatalf("Failed to bind UDP: %s",err)
	}

	// send the request for info
	n, err := conn.Write([]byte{Req})
	if err != nil {
		log.Fatalf("Failed to send request: %s",err)
	}
	log.Printf("Wrote %d bytes to %s",n,rem)

	// server response - 1 byte for type, 4 for the port, and 16 for the IP
	resp := make([]byte, 25)
	n, err = conn.Read(resp)
	if err != nil {
		log.Fatalf("Failed to read response: %s",err)
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
		log.Fatal("Failed to read integer from portBytes")
	}

	return &net.UDPAddr{ IP: ipBytes, Port: int(extPort) }


}




