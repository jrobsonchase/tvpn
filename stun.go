package tvpn

import "net"

type StunBackend interface {
	DiscoverExt(port int) (net.IP,int)
}
