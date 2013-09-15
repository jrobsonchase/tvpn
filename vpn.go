package tvpn


type VPNBackend interface {
	Connect(remote,localtun string,remoteport,localport int, key [][]byte, dir bool) VPNConn
}

type VPNConn interface {
	Disconnect()
	Status() int
}
