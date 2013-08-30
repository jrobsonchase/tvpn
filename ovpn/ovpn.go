package ovpn

import (
	"io"
	"log"
	"os/exec"
	"os"
)

type OVPN struct {
	Path		string
	RemoteIP	string
	RemotePort	string
	LocalPort	string
	RemoteTunIP	string
	LocalTunIP	string
	KeyFile		string
	Direction	bool
	LogFile		string
	ErrFile		string
	Cmd			*exec.Cmd
}

var ovpnOpts []string = []string{
	"--mode","p2p",
	"--proto","udp",
	"--dev","tap",
	"--ping-exit","30",
	"--ping","10"}

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
		"--remote",v.RemoteIP,
		"--rport",v.RemotePort,
		"--lport",v.LocalPort,
		"--secret",v.KeyFile,dir
		"--ifconfig",v.LocalTunIP,"255.255.255.252")...)

	e := cmd.Run()

	if e != nil {
		log.Fatal(e.Error())
	}
	logFile,err := os.OpenFile(v.LogFile,os.O_APPEND,0666)
	if err != nil {
		log.Fatalf("Failed to open file for writing: %s",err)
	}
	errFile,err := os.OpenFile(v.ErrFile,os.O_APPEND,0666)
	if err != nil {
		log.Fatalf("Failed to open file for writing: %s",err)
	}

	go io.Copy(cmd.Stdout,logFile)
	go io.Copy(cmd.Stderr,errFile)

	v.Cmd = cmd
	return
}

func (v *OVPN) Disconnect() {
	return
}

func (v *OVPN) Restart() {
	return
}
