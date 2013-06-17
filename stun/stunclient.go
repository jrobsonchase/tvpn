package main

import (
	"net"
	"log"
	"tvpn/stun"
)

func main() {
	remote,_ := net.ResolveUDPAddr("udp","hyperion.chthonius.net:12345")
	resp := stun.DiscoverExternal(22222,remote)

	log.Printf("recieved ip %s port %d",resp.IP.String(), resp.Port)
}
