package ovpn

import (
	"io"
	"log"
	. "net"
	"os"
	"os/exec"
)

type oVPN struct {
	path		string
	remoteAddr	TCPAddr
	certFile	string
	keyFile		string
	proc		*os.Process
}

func spawnWriteChan(c io.WriteCloser) (o chan string) {
	return
}

func spawnReadChan(c io.ReadCloser) (o chan string) {
	return
}

func spawnReadWriteChan(c io.ReadWriteCloser) (o chan string) {
	return
}

func (v *oVPN) Connect() (in,out,err chan string) {
	cmd := exec.Command(v.path)
	e := cmd.Run()
	if e != nil {
		log.Fatal(e.Error())
	}
	v.proc = cmd.Process
	in = make(chan string)
	out = make(chan string)
	err = make(chan string)
	return
}

func (oVPN) Disconnect() {
	return
}

func (oVPN) Restart() {
	return
}

func New(addr,port,cert,key string) oVPN {
	return oVPN{}
}
