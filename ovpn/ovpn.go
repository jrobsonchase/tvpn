package ovpn

import (
	"fmt"
	"io"
	"runtime"
	"log"
	"os"
	"os/exec"
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

func New() OVPNBackend {
	var ovpnpath string
	var tmp string
	if runtime.GOOS == "windows" {
		tmp = "%TEMP%\\"
		ovpnpath = "ovpn/openvpn.exe"
	} else {
		tmp = "/tmp/"
		ovpnpath = "/usr/bin/openvpn"
	}
	return OVPNBackend{ovpnpath,tmp}
}

func (ovpn *OVPNBackend) Connect(remoteip,localtun string,
	remoteport,localport int,
	key [][]byte,
	dir bool) (*OVPNConn,error) {

	var dirS string
	if dir {
		dirS = "1"
	} else {
		dirS = "0"
	}

	keyfile := fmt.Sprintf("%s%s.key",ovpn.tmp,remoteip)
	keyhandle,err := os.Create(keyfile)
	if err != nil {
		return nil,err
	}
	_,err = keyhandle.Write(EncodeOpenVPNKey(key))
	if err != nil {
		return nil,err
	}
	keyhandle.Close()


	cmd := exec.Command(ovpn.path,
		append(ovpnOpts,
			"--remote", remoteip,
			"--rport", fmt.Sprintf("%d",remoteport),
			"--lport", fmt.Sprintf("%d",localport),
			"--secret", keyfile, dirS,
			"--ifconfig", localtun, "255.255.255.252")...)

	e := cmd.Run()

	if e != nil {
		log.Fatal(e.Error())
	}

	logFile, err := os.OpenFile(fmt.Sprintf("%s%s.err",ovpn.tmp,remoteip), os.O_APPEND, 0666)
	if err != nil {
		return nil,err
	}
	errFile, err := os.OpenFile(fmt.Sprintf("%s%s.err",ovpn.tmp,remoteip), os.O_APPEND, 0666)
	if err != nil {
		return nil,err
	}

	go io.Copy(cmd.Stdout, logFile)
	go io.Copy(cmd.Stderr, errFile)


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
