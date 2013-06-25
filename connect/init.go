package connect

import (
	"time"
	"net"
	"math/rand"
)

var StunServer *net.UDPAddr
var rgen *rand.Rand
var PubKey string

func init() {
	rgen = rand.New(rand.NewSource(time.Now().UnixNano()))
}


