package main

import (
	"tvpn/stun"
)

func main() {
	stun.ServeStun(12345)
}
