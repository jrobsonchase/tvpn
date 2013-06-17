package stun

import (
	"net"
	"log"
	"encoding/binary"
)

func ServeStun(port int) {
	network := "udp"
	laddr := &net.UDPAddr{ Port: port }
	conn, err := net.ListenUDP(network,laddr)
	if err != nil {
		log.Fatalf("Failed to bind UDP port: %s",err)
	}

	req := make([]byte,1)
	portBytes := make([]byte,8)
	for {
		n,raddr,err := conn.ReadFromUDP(req)
		if err != nil {
			log.Fatalf("Failed to read from remote: %s",err)
		}

		n = binary.PutVarint(portBytes,int64(raddr.Port))
		if n == 0 {
			log.Fatal("Failed to convert port to bytes")
		}

		resp := append([]byte{Resp},portBytes...)
		resp = append(resp,raddr.IP...)
		_,err = conn.WriteTo(resp,raddr)
		if err != nil {
			log.Fatalf("Failed to send response: %s",err)
		}
		log.Printf("Sent response to %s at port %d",raddr.IP.String(),raddr.Port)
	}
}

