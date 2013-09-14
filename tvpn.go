package tvpn

import (
	"fmt"
	"net"
)

type TVPN struct {
	Name, IRCServer, STUNServer, Group string
	Signaling                          Backend
	Friends                            []string
}

type Conn net.Conn

func (t TVPN) TestStun() error {
	return nil
}

func (t TVPN) Connect() (*Conn, error) {
	return nil, nil
}

func (t TVPN) Run() error {
	fmt.Printf("Hello, World!\n")
	return nil
}
