package tvpn

import (
	"github.com/Pursuit92/LeveledLogger/log"
	"math/rand"
	"time"
)

type TVPN struct {
	Name, Group string
	Friends     []string
	Sig         SigBackend
	Stun        StunBackend
	VPN			VPNBackend
	States      map[string]*ConState
	Alloc		IPManager
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

func New(name, group string,
	friends []string,
	sig SigBackend,
	stun StunBackend,
	vpn VPNBackend) *TVPN {

	return nil
}

func (t *TVPN) Run() error {
	for {
		msg := t.Sig.RecvMessage()
		switch msg.Type {
		case Init:
			t.States[msg.From] = &ConState{}
			friend := false
			for _,v := range t.Friends {
				if v == msg.From {
					friend = true
				}
			}
			t.States[msg.From].InitState(msg.From,friend,false,t.Sig)
		case Join:
			for _,v := range t.Friends {
				if v == msg.From {
					t.States[msg.From] = &ConState{}
					t.States[msg.From].InitState(msg.From,true,true,t.Sig)
				}
			}

		case Quit:
			delete(t.States,msg.From)
		case Reset:
			t.States[msg.From].Reset(t.Sig,msg.Data["reason"])
		default:
			for i,v := range t.States {
				if i == msg.From {
					v.Input(msg,*t)
				}
			}
		}
	}
}
