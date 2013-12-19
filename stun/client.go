package stun

import (
	"net"
	"fmt"
	s "github.com/Pursuit92/stun"
	"github.com/Pursuit92/LeveledLogger/log"
)

type StunErr string

func (s StunErr) Error() string {
	return string(s)
}

type StunBackend struct {
	Server string
}

func (b StunBackend) DiscoverExt(port int) (net.IP,int,error) {
	uadd := DiscoverExternal(port,b.Server)
	if uadd == nil {
		return nil,0,StunErr("Failed to get external address!")
	}
	return uadd.IP,uadd.Port,nil
}

func DiscoverExternal(port int, addr string) (*net.UDPAddr) {
	message := s.NewMessage()

	message.Class = s.Request
	message.Method = s.Binding

	message.AddAttribute(s.MappedAddress("0.0.0.0",port))

	resp := s.SendMessage(message,fmt.Sprintf("0.0.0.0:%d",port),addr)

	for _,v := range resp.Attrs {
		if v.Type == s.MappedAddressCode {
			ma,ok := v.Attr.(s.MappedAddressAttr)
			if ok {
				return &net.UDPAddr{ IP: ma.Address, Port: ma.Port }
			}
		}
	}

	return nil
}

func SetLogLevel(n int) {
	log.Out.SetLevel(n)
}

func SetLogPrefix(s string) {
	log.Out.SetPrefix(2,s)
	log.Out.SetPrefix(3,s)
}
