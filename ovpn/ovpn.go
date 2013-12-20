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

package ovpn

import (
	"fmt"
	"math/big"
	"io"
	"runtime"
	"log"
	"os"
	"os/exec"
	"github.com/Pursuit92/tvpn"
)

type OVPNBackend struct {
	path,tmp string
}

/*
	Path        string
	RemoteIP    string
	RemotePort  string
	LocalPort   string
	RemoteTunIP string
	LocalTunIP  string
	KeyFile     string
	LogFile     string
	ErrFile     string
	Cmd         *exec.Cmd
}
*/

type OVPNConn struct {
	Cmd *exec.Cmd
}

func New() *OVPNBackend {
	var ovpnpath string
	var tmp string
	if runtime.GOOS == "windows" {
		//tmp = "%TEMP%\\"
		tmp = os.ExpandEnv("${TEMP}\\")
		ovpnpath = `C:\Program Files (x86)\OpenVPN\bin\openvpn.exe`
	} else {
		tmp = "/tmp/"
		ovpnpath = "/usr/bin/openvpn"
	}
	return &OVPNBackend{ovpnpath,tmp}
}

func (ovpn *OVPNBackend) Connect(remoteip,localtun string,
	remoteport,localport int,
	key []*big.Int,
	dir bool) (tvpn.VPNConn,error) {

	var dirS string
	if dir {
		dirS = "1"
	} else {
		dirS = "0"
	}

	keyfile := fmt.Sprintf("%s%s-%d.key",ovpn.tmp,remoteip,remoteport)
	keyhandle,e := os.Create(keyfile)
	if e != nil {
		log.Fatal(e)
		return nil,e
	}
	_,e = keyhandle.Write(EncodeOpenVPNKey(key))
	if e != nil {
		log.Fatal(e)
		return nil,e
	}
	keyhandle.Close()


	opts := append(ovpnOpts,
			"--remote", remoteip,
			"--rport", fmt.Sprintf("%d",remoteport),
			"--lport", fmt.Sprintf("%d",localport),
			"--secret", keyfile, dirS,
			"--ifconfig", localtun, "255.255.255.252")

	cmd := exec.Command(ovpn.path, opts...)


	fmt.Printf("Running command: %s ",cmd.Path)
	for _,v := range cmd.Args {
		fmt.Printf("%s ",v)
	}
	fmt.Print("\n")
	out,e := cmd.StdoutPipe()

	if e != nil {
		log.Fatal(e)
	}
	err,e := cmd.StderrPipe()

	if e != nil {
		log.Fatal(e)
	}
	e = cmd.Start()
	if e != nil {
		log.Fatal(e)
	}


	/*
	logFile, e := os.Create(fmt.Sprintf("%s%s.log",ovpn.tmp,remoteip))
	if e != nil {
		log.Fatal(e)
		return nil,e
	}
	errFile, e := os.Create(fmt.Sprintf("%s%s.err",ovpn.tmp,remoteip))
	if e != nil {
		log.Fatal(e)
		return nil,e
	}
	*/

	go io.Copy(os.Stdout,out)
	go io.Copy(os.Stderr,err)

	log.Printf("\nVPN Connected with pid %d\n",cmd.Process.Pid)
	return &OVPNConn{cmd},nil
}

func (conn *OVPNConn) Disconnect() {
	proc := conn.Cmd.Process
	proc.Kill()
}

func (conn OVPNConn) Status() int {
	return 0
}

var ovpnOpts []string = []string{
	"--mode", "p2p",
	"--proto", "udp",
	"--dev", "tap",
	"--ping-exit", "30",
	"--ping", "10",
}

/*
func (v *OVPN) Connect() {
	// Set direction for --secret option - allows all of the secret to be used
	var dir string
	if v.Direction {
		dir = "1"
	} else {
		dir = "0"
	}

	cmd := exec.Command(v.Path,
		append(ovpnOpts,
			"--remote", v.RemoteIP,
			"--rport", v.RemotePort,
			"--lport", v.LocalPort,
			"--secret", v.KeyFile, dir,
			"--ifconfig", v.LocalTunIP, "255.255.255.252")...)

	e := cmd.Run()

	if e != nil {
		log.Fatal(e.Error())
	}
	logFile, err := os.OpenFile(v.LogFile, os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open file for writing: %s", err)
	}
	errFile, err := os.OpenFile(v.ErrFile, os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open file for writing: %s", err)
	}

	go io.Copy(cmd.Stdout, logFile)
	go io.Copy(cmd.Stderr, errFile)

	v.Cmd = cmd
	return
}

func (v *OVPN) Disconnect() {
	return
}

func (v *OVPN) Restart() {
	return
}
*/
