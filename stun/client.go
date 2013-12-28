/*
 *  TVPN: A Peer-to-Peer VPN solution for traversing NAT firewalls
 *  Copyright (C) 2013  Joshua Chase <jcjoshuachase@gmail.com>
 *
 *  This program is free software; you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation; either version 2 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License along
 *  with this program; if not, write to the Free Software Foundation, Inc.,
 *  51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/

package stun

import (
	"net"
	"fmt"
	s "github.com/Pursuit92/stun"
	"github.com/Pursuit92/tvpn"
	"github.com/Pursuit92/LeveledLogger/log"
)

type StunErr string

func (s StunErr) Error() string {
	return string(s)
}

type StunBackend struct {
	Server string
}

func (b *StunBackend) Configure(conf tvpn.StunConfig) {
	b.Server = conf["Server"]
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
