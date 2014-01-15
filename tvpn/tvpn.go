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

package main

import (
	"flag"
	"fmt"
	"os"
	//"runtime"
	"github.com/Pursuit92/tvpn"
	"github.com/Pursuit92/tvpn/ovpn"
	"github.com/Pursuit92/tvpn/stun"
	"github.com/Pursuit92/tvpn/irc"
	"github.com/Pursuit92/LeveledLogger/log"
)

const friendLimit int = 256

func exitError(s string) {
	fmt.Printf("%s\n", s)
	os.Exit(1)
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	verboseLevel := flag.Int("v",1,"Verbosity level. Set to 1 by default")
	configPath := flag.String("config","/usr/share/tvpn/tvpn.config","JSON Configuration file")
	flag.Parse()

	conf,err := tvpn.ReadConfig(*configPath)
	if err != nil {
		exitError(err.Error())
	}


	friends := make([]string,len(conf.Friends))
	i := 0
	for f,_ := range conf.Friends {
		friends[i] = f
		i++
	}

	switch *verboseLevel {
	case 0,1:
	case 2:
		irc.SetLogLevel(2)
		tvpn.SetLogLevel(2)
	case 3:
		stun.SetLogLevel(2)
		irc.SetLogLevel(2)
		tvpn.SetLogLevel(2)
	case 4:
		stun.SetLogLevel(2)
		irc.SetLogLevel(2)
		tvpn.SetLogLevel(3)
	case 5:
		stun.SetLogLevel(3)
		irc.SetLogLevel(3)
		tvpn.SetLogLevel(3)
	case 6,7,8,9,10:
		stun.SetLogLevel(10)
		irc.SetLogLevel(10)
		tvpn.SetLogLevel(10)
	default:
	}

	log.Out.Printf(0,"Loaded friends:\n")
	for _, v := range friends {
		log.Out.Println(0,v)
	}

	conf.Sig["Group"] = conf.Group
	conf.Sig["Name"] = conf.Name

	ipman := new(tvpn.IPManager)
	stunman := new(stun.StunBackend)
	vpnman := new(ovpn.OVPNBackend)
	ircman := new(irc.IRCBackend)

	/*
	ipman.Configure(conf.IPMan)
	stunman.Configure(conf.Stun)
	vpnman.Configure(conf.VPN)
	ircman.Configure(conf.Sig)
	*/

	if err != nil {
		exitError(err.Error())
	}

	tvpnInstance := tvpn.New(
		ircman,
		stunman,
		vpnman,
		ipman)

	tvpnInstance.Configure(*conf)

	err = tvpnInstance.Run()

}
