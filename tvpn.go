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
	"math/rand"
	"time"
	"fmt"
)

type Friend struct {
	Validate bool
	Routes map[string]string
}

type TVPN struct {
	Friends     map[string]Friend
	Sig         SigBackend
	Stun        StunBackend
	VPN			VPNBackend
	States      map[string]*ConState
	Alloc		*IPManager
}

var rgen *rand.Rand
func init() {
	rgen = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func SetLogLevel(n int) {
	log.Out.SetLevel(n)
}

func SetLogPrefix(s string) {
	for i := 0; i < 10; i++ {
		log.Out.SetPrefix(i,s)
	}
}

func New(sig SigBackend, stun StunBackend, vpn VPNBackend, alloc *IPManager) *TVPN {
	tvpnInstance := TVPN{
		Sig:  sig,
		Stun: stun,
		VPN: vpn,
		Alloc: alloc,
		States: make(map[string]*ConState),
	}

	return &tvpnInstance
}

func (t *TVPN) Configure(conf Config) {
	t.Friends = conf.Friends
	t.Sig.Configure(conf.Sig)
	t.Stun.Configure(conf.Stun)
	t.VPN.Configure(conf.VPN)
	t.Alloc.Configure(conf.IPMan)
}

func (t TVPN) IsFriend(name string) (Friend,bool) {
	f,ok := t.Friends[name]
	return f,ok
}


func (t *TVPN) Run() error {
	for {
		fmt.Printf("Waiting for message...\n")
		msg := t.Sig.RecvMessage()
		fmt.Printf("Got a message: %s\n",msg.String())
		switch msg.Type {
		case Init:
			friend,ok := t.IsFriend(msg.From)
			fmt.Printf("Creating new state machine for %s\n",msg.From)
			t.States[msg.From] = NewState(msg.From,friend,ok,false,*t)
			t.States[msg.From].Input(msg,*t)
		case Join:
			fmt.Printf("Received Join from %s!\n",msg.From)
			friend,ok := t.IsFriend(msg.From)
			if ok {
				t.States[msg.From] = NewState(msg.From,friend,true,true,*t)
			}
			fmt.Println("Done with join!")

		case Quit:
			st,exists := t.States[msg.From]
			if exists {
				st.Cleanup(*t)
				delete(t.States,msg.From)
			}
		case Reset:
			t.States[msg.From].Reset(msg.Data["reason"],*t)
		default:
			st,exists := t.States[msg.From]
			if exists {
				st.Input(msg,*t)
			} else {
				// do stuff here
			}
		}
	}
}
