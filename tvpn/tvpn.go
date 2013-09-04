package tvpn

import "net"
type TVPN struct {
	Name,IRCServer,STUNServer,Group string
	Friends []string
}

type Conn net.Conn

func (t TVPN) TestStun() error {
	return nil
}

func (t TVPN) Connect() (*Conn,error) {
	return nil,nil
}
