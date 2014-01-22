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

package tvpn

import (
	"github.com/Pursuit92/LeveledLogger/log"
	"fmt"
	"net"
	"github.com/Pursuit92/tvpn/dh"
	"github.com/Pursuit92/state"
)

type ConState struct {
	Name         string
	Init         bool
	tvpn		 *TVPN
	Params       []dh.Params
	Key          [][64]byte
	keys		 int
	Tun          net.IP
	Port		 int
	Conn		 VPNConn
	mach		 state.Simple
	resets		 int
}

func (c *ConState) Push(s ...state.StateTrans) {
	c.mach.Push(s...)
}

func (c *ConState) Pop() (state.StateTrans,error) {
	return c.mach.Pop()
}

func (st *ConState) Reset(reason string) {
	st.Cleanup()
	if reason != "" {
		log.Out.Lprintf(3,"Conversation with %s reset. Reason: %s\n",st.Name,reason)
	}
	st.resets++
	if st.resets > 3 {
		st.Push(trapState)
	} else {
		*st = *(NewState(st.Name,st.Init,st.tvpn))
	}
}

func NewState(name string, init bool,t *TVPN) *ConState {
	st := &ConState{Name: name,
		Init: init,
		tvpn: t}
	if init {
		t.Sig.SendMessage(Message{Type: Init, To: name})
		st.Push(initState)
	} else {
		st.Push(noneState)
	}
	return st
}

func trapState(sm state.StateMachine,mesInt interface{}) error {
	sm.Push(trapState)
	return nil
}

// NoneState is the state in which we wait for an Init
// Next state is DHNeg after a valid Init
func noneState(sm state.StateMachine,mesInt interface{}) error {
	st := sm.(*ConState)
	t := st.tvpn
	mes := mesInt.(Message)
	switch mes.Type {
	case Init:
		_, ok := t.IsFriend(st.Name)
		if ok {
			t.Sig.SendMessage(Message{Type: Accept, To: st.Name})
			st.Push(dhNegState)
			st.Params = make([]dh.Params, 4)
			st.Key = make([][64]byte, 4)
			for i := 0; i < 4; i++ {
				st.Params[i] = dh.GenParams()
				t.Sig.SendMessage(Message{To: st.Name, Type: Dhpub, Data: map[string]string{
					"i": fmt.Sprintf("%d", i),
					"x": st.Params[i].XS(),
					"y": st.Params[i].YS(),
				}})
			}
		} else {
			t.Sig.SendMessage(Message{To: st.Name, Type: Deny, Data: map[string]string{"reason": "Not Authorized"}})
			st.Push(noneState)
		}
	default:
		t.Sig.SendMessage(Message{To: st.Name, Type: Reset, Data: map[string]string{
			"reason": "Invalid state: None"}})
			st.Push(noneState)
	}
	return nil
}

// Init state is after Init is sent and before Accept is received
// Next state is DHNeg
func initState(sm state.StateMachine,mesInt interface{}) error {
	st := sm.(*ConState)
	t := st.tvpn
	mes := mesInt.(Message)
	switch mes.Type {
	case Accept:
		st.Push(dhNegState)
		st.Params = make([]dh.Params, 4)
		st.Key = make([][64]byte, 4)
		for i := 0; i < 4; i++ {
			st.Params[i] = dh.GenParams()
			t.Sig.SendMessage(Message{To: st.Name, Type: Dhpub, Data: map[string]string{
				"i": fmt.Sprintf("%d", i),
				"x": st.Params[i].XS(),
				"y": st.Params[i].YS(),
			}})
		}
	case Deny:
		st.Push(trapState)
	default:
		t.Sig.SendMessage(Message{To: st.Name, Type: Reset, Data: map[string]string{
			"reason": "Invalid state: Init"}})
			st.Reset("")
	}
	return nil
}

func dhNegState(sm state.StateMachine,mesInt interface{}) error {
	st := sm.(*ConState)
	t := st.tvpn
	mes := mesInt.(Message)
	switch mes.Type {
	case Dhpub:
		x, y, i, err := mes.DhParams()
		if err != nil {
			t.Sig.SendMessage(Message{To: st.Name, Type: Reset, Data: map[string]string{
				"reason": "Invalid DH Params",
			}})
			st.Reset("")
			return nil
		}
		st.Key[i] = dh.GenKey(st.Params[i], dh.Params{X: x, Y: y})
		st.keys++
		if st.keys < 4 {
			st.Push(dhNegState)
			return nil
		}
		st.Tun = t.Alloc.Request(nil)
		t.Sig.SendMessage(Message{Type: Tunnip, To: st.Name, Data: map[string]string{"ip": st.Tun.String()}})
		st.Push(tunNegState)
	default:
		t.Sig.SendMessage(Message{Type: Reset, Data: map[string]string{
			"reason": "Invalid state: DHNeg"}})
			st.Reset("")

	}
	return nil
}

func tunNegState(sm state.StateMachine,mesInt interface{}) error {
	st := sm.(*ConState)
	t := st.tvpn
	mes := mesInt.(Message)
	switch mes.Type {
	case Tunnip:
		ip,_ := mes.IPInfo()
		if ! ip.Equal(st.Tun) {
			if isGreater(ip,st.Tun) {
				t.Alloc.Release(st.Tun)
				st.Tun = t.Alloc.Request(ip)
			}
			t.Sig.SendMessage(Message{Type: Tunnip, To: st.Name, Data: map[string]string{"ip": st.Tun.String()}})
			st.Push(tunNegState)
			return nil
		}
		st.Port = rgen.Int() % (65536 - 49152) + 49152
		ip,port,err := t.Stun.DiscoverExt(st.Port)
		if err != nil {
			t.Sig.SendMessage(Message{Type: Reset, Data: map[string]string{
				"reason": "Failed to discover external connection info"}})
			st.Reset("")
			return nil
		}
		t.Sig.SendMessage(Message{Type: Conninfo, To: st.Name, Data: map[string]string{
			"port": fmt.Sprintf("%d",port),
			"ip": ip.String(),
		}})
		st.Push(conNegState)
	default:
		t.Sig.SendMessage(Message{Type: Reset, Data: map[string]string{
			"reason": "Invalid state: DHNeg"}})
			st.Reset("")
	}
	return nil
}

func conNegState(sm state.StateMachine,mesInt interface{}) error {
	st := sm.(*ConState)
	t := st.tvpn
	mes := mesInt.(Message)
	switch mes.Type {
	case Conninfo:
		ip,port := mes.IPInfo()
		log.Out.Lprintf(2,"Connecting vpn...")
		friend, _ := t.IsFriend(st.Name)
		conn, err := t.VPN.Connect(ip,st.Tun,port,st.Port,st.Key,st.Init,friend.Routes)
		if err == nil {
			log.Out.Lprintf(2,"VPN Connected!\n")
			st.Conn = conn
			st.Push(connectedState)
		} else {
			log.Out.Lprintf(2,"Error connecting VPN: %s\n",err.Error())
			st.Reset("")
		}

	default:
		t.Sig.SendMessage(Message{Type: Reset, Data: map[string]string{
			"reason": "Invalid state: DHNeg"}})
			st.Reset("")
	}
	return nil
}


func connectedState(sm state.StateMachine,mesInt interface{}) error {
	sm.Push(connectedState)
	return nil
}

func (st *ConState) Cleanup() {
	t := st.tvpn
	if st.Conn != nil {
		st.Conn.Disconnect()
		st.Conn = nil
	}

	if st.Tun != nil {
		t.Alloc.Release(st.Tun)
		st.Tun = nil
	}
}
