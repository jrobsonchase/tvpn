package ovpn

import (
	"io"
	"log"
	"os/exec"
)

type OVPN struct {
	Path		string
	RemoteIP	string
	RemotePort	string
	LocalPort	string
	RemoteTunIP	string
	LocalTunIP	string
	CertFile	string
	KeyFile		string
	LogFile		string
	ErrFile		string
	Cmd			*exec.Cmd
}

// Functions to spawn jobs that turn readers/writers to channels
func spawnWriteChan(c io.WriteCloser) (o chan string) {
	return
}

func spawnReadChan(c io.ReadCloser) (o chan string) {
	return
}

func spawnReadWriteChan(c io.ReadWriteCloser) (o chan string) {
	return
}

func (v *OVPN) Connect() {
	cmd := exec.Command(v.Path,
		"--mode","p2p",
		"--proto","udp",
		"--dev-type","tap",
		"--remote",v.RemoteIP,
		"--rport",v.RemotePort,
		"--lport",v.LocalPort,
		"--cert",v.CertFile,
		"--key",v.KeyFile,
		"--ifconfig",v.LocalTunIP,"255.255.255.252",
		"--route",v.RemoteIP,"255.255.255.255",v.RemoteTunIP)

	e := cmd.Run()

	if e != nil {
		log.Fatal(e.Error())
	}

	v.Cmd = cmd
	return
}

func (v *OVPN) Disconnect() {
	return
}

func (v *OVPN) Restart() {
	return
}
