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
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"github.com/Pursuit92/tvpn/dh"
)

type ConState struct {
	State        int
	Name         string
	Params       []dh.Params
	Key          []*big.Int
	IP           net.IP
	Port		 int
	Friend, Init bool
	Data		 Friend
	Conn		 VPNConn
}

const (
	NoneState int = iota
	InitState
	DHNeg
	TunNeg
	ConNeg
	Connected
	DeleteMe
)

func (st *ConState) Input(mes Message, t TVPN) {
	fmt.Printf("Got message: %s\n",mes.String())
	switch st.State {
	case NoneState:
		fmt.Printf("in NoneState\n")
		st.noneState(mes,t)
		fmt.Printf("Done with state update!\n")
	case InitState:
		fmt.Printf("in InitState\n")
		st.initState(mes,t)
		fmt.Printf("Done with state update!\n")
	case DHNeg:
		fmt.Printf("in DHNeg\n")
		st.dhnegState(mes,t)
		fmt.Printf("Done with state update!\n")
	case TunNeg:
		fmt.Printf("in TunNeg\n")
		st.tunnegState(mes,t)
		fmt.Printf("Done with state update!\n")
	case ConNeg:
		fmt.Printf("in ConNeg\n")
		st.connegState(mes,t)
		fmt.Printf("Done with state update!\n")
	case Connected:
		fmt.Printf("in Connected\n")
		st.connectedState(mes,t)
		fmt.Printf("Done with state update!\n")
	default:
	}
}

func (st *ConState) Reset(reason string, t TVPN) {
	st.Cleanup(t)
	*st = *(NewState(st.Name,st.Data,st.Friend,st.Init,t))
	if reason != "" {
		log.Out.Printf(3,"Conversation with %s reset. Reason: %s\n",st.Name,reason)
	}
}

func NewState(name string, fData Friend, friend,init bool,t TVPN) *ConState {
	st := ConState{}
	st.Name = name
	st.Friend = friend
	st.Init = init
	st.Data = fData
	if init {
		t.Sig.SendMessage(Message{Type: Init, To: name})
		st.State = InitState
	} else {
		st.State = NoneState
	}
	return &st
}

// NoneState is the state in which we wait for an Init
// Next state is DHNeg after a valid Init
func (st *ConState) noneState(mes Message, t TVPN) {
	switch mes.Type {
	case Init:
		if st.Friend {
			t.Sig.SendMessage(Message{Type: Accept, To: st.Name})
			st.State = DHNeg
			st.Params = make([]dh.Params, 4)
			st.Key = make([]*big.Int, 4)
			for i := 0; i < 4; i++ {
				st.Params[i] = dh.GenParams()
				t.Sig.SendMessage(Message{To: st.Name, Type: Dhpub, Data: map[string]string{
					"i": fmt.Sprintf("%d", i),
					"x": base64.StdEncoding.EncodeToString(st.Params[i].X.Bytes()),
					"y": base64.StdEncoding.EncodeToString(st.Params[i].Y.Bytes()),
				}})
			}
		} else {
			t.Sig.SendMessage(Message{To: st.Name, Type: Deny, Data: map[string]string{"reason": "Not Authorized"}})
			st.State = DeleteMe
		}
	default:
		t.Sig.SendMessage(Message{To: st.Name, Type: Reset, Data: map[string]string{
			"reason": "Invalid state: None"}})
			st.Reset("",t)
	}
}

// Init state is after Init is sent and before Accept is received
// Next state is DHNeg
func (st *ConState) initState(mes Message, t TVPN) {
	switch mes.Type {
	case Accept:
		st.State = DHNeg
		st.Params = make([]dh.Params, 4)
		st.Key = make([]*big.Int, 4)
		for i := 0; i < 4; i++ {
			st.Params[i] = dh.GenParams()
			t.Sig.SendMessage(Message{To: st.Name, Type: Dhpub, Data: map[string]string{
				"i": fmt.Sprintf("%d", i),
				"x": base64.StdEncoding.EncodeToString(st.Params[i].X.Bytes()),
				"y": base64.StdEncoding.EncodeToString(st.Params[i].Y.Bytes()),
			}})
		}
	default:
		t.Sig.SendMessage(Message{To: st.Name, Type: Reset, Data: map[string]string{
			"reason": "Invalid state: Init"}})
			st.Reset("",t)
	}
}

func (st *ConState) dhnegState(mes Message, t TVPN) {
	switch mes.Type {
	case Dhpub:
		x, y, i, err := mes.DhParams()
		if err != nil {
			t.Sig.SendMessage(Message{To: st.Name, Type: Reset, Data: map[string]string{
				"reason": "Invalid DH Params",
			}})
			st.Reset("",t)
			return
		}
		st.Key[i] = dh.GenMutSecret(st.Params[i], dh.Params{X: x, Y: y})
		for _, v := range st.Key {
			if v == nil {
				// end state change - still need more keys
				return
			}
		}
		st.IP = t.Alloc.Request(nil)
		t.Sig.SendMessage(Message{Type: Tunnip, To: st.Name, Data: map[string]string{"ip": st.IP.String()}})
		st.State = TunNeg
	default:
		t.Sig.SendMessage(Message{Type: Reset, Data: map[string]string{
			"reason": "Invalid state: DHNeg"}})
			st.Reset("",t)

	}
}

func (st *ConState) tunnegState(mes Message, t TVPN) {
	switch mes.Type {
	case Tunnip:
		ip,_ := mes.IPInfo()
		if ! ip.Equal(st.IP) {
			fmt.Printf("IP's not equal!\n")
			if isGreater(ip,st.IP) {
				t.Alloc.Release(st.IP)
				st.IP = t.Alloc.Request(ip)
			}
			t.Sig.SendMessage(Message{Type: Tunnip, To: st.Name, Data: map[string]string{"ip": st.IP.String()}})
			return
		}
		st.Port = rgen.Int() % (65536 - 49152) + 49152
		ip,port,err := t.Stun.DiscoverExt(st.Port)
		if err != nil {
			t.Sig.SendMessage(Message{Type: Reset, Data: map[string]string{
				"reason": "Failed to discover external connection info"}})
			st.Reset("",t)
			return
		}
		t.Sig.SendMessage(Message{Type: Conninfo, To: st.Name, Data: map[string]string{
			"port": fmt.Sprintf("%d",port),
			"ip": ip.String(),
		}})
		st.State = ConNeg
	default:
		t.Sig.SendMessage(Message{Type: Reset, Data: map[string]string{
			"reason": "Invalid state: DHNeg"}})
			st.Reset("",t)
	}

}

func (st *ConState) connegState(mes Message,t TVPN) {
	switch mes.Type {
	case Conninfo:
		ip,port := mes.IPInfo()
		fmt.Printf("Connecting vpn...")
		conn, err := t.VPN.Connect(ip,st.IP,port,st.Port,st.Key,st.Init,st.Data.Routes)
		if err == nil {
			fmt.Printf("VPN Connected!\n")
			st.Conn = conn
			st.State = Connected
		} else {
			fmt.Printf("Error connecting VPN: %s\n",err.Error())
		}

	default:
		t.Sig.SendMessage(Message{Type: Reset, Data: map[string]string{
			"reason": "Invalid state: DHNeg"}})
			st.Reset("",t)
	}
}


func (st *ConState) connectedState(mes Message,t TVPN) {
}

func (st *ConState) Cleanup(t TVPN) {
	if st.Conn != nil {
		st.Conn.Disconnect()
		st.Conn = nil
	}

	if st.IP != nil {
		t.Alloc.Release(st.IP)
		st.IP = nil
	}
}
