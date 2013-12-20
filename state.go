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
	*st = *(NewState(st.Name,st.Friend,st.Init,t))
	if reason != "" {
		log.Out.Printf(3,"Conversation with %s reset. Reason: %s\n",st.Name,reason)
	}
}

func NewState(name string,friend,init bool,t TVPN) *ConState {
	st := ConState{}
	st.Name = name
	st.Friend = friend
	st.Init = init
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
		if ! isEqual(ip,st.IP) {
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
		tunIP := st.IP
		if st.Init {
			tunIP[len(tunIP)-1] += 1
		} else {
			tunIP[len(tunIP)-1] += 2
		}
		conn, err := t.VPN.Connect(ip.String(),tunIP.String(),port,st.Port,st.Key,st.Init)
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
