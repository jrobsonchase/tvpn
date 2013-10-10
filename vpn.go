package tvpn

import "math/big"

type VPNBackend interface {
	Connect(remote,localtun string,remoteport,localport int, key []*big.Int, dir bool) VPNConn
}

type VPNConn interface {
	Disconnect()
	Status() int
}
