package main

import "tvpn/ovpn"

func main() {
	print(&ovpn.OVPN{RemoteIP: "1.2.3.4"})
}
