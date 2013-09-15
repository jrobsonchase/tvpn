package tvpn

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"tvpn/ovpn"
)

type ConState struct {
	State        int
	Name         string
	Params       []ovpn.Params
	Key          []*big.Int
	IP           net.IP
	Friend, Init bool
}

const (
	NoneState int = iota
	InitState
	DHNeg
	TunNeg
	Connected
	Error
)

func (st *ConState) Input(mes Message, tvpn TVPN) {
	switch st.State {
	case NoneState:
		st.noneState(mes, tvpn)
	case InitState:
		st.initState(mes, tvpn)
	case DHNeg:
		st.dhnegState(mes, tvpn)
	case TunNeg:
		st.tunnegState(mes, tvpn)
	case Connected:
		st.connectedState(mes, tvpn)
	}
}

// NoneState is the state in which we wait for an Init
// Next state is DHNeg after a valid Init
func (st *ConState) noneState(mes Message, tvpn TVPN) {
	switch mes.Type {
	case Init:
		if st.Friend {
			tvpn.Sig.SendMessage(Message{Type: Accept, To: st.Name})
			st.Params = make([]ovpn.Params, 4)
			st.Key = make([]*big.Int, 4)
			for i := 0; i < 4; i++ {
				st.Params[i] = ovpn.GenParams()
				tvpn.Sig.SendMessage(Message{Type: Dhpub, Data: map[string]string{
					"i": fmt.Sprintf("%d", i),
					"x": base64.StdEncoding.EncodeToString(st.Params[i].X.Bytes()),
					"y": base64.StdEncoding.EncodeToString(st.Params[i].Y.Bytes()),
				}})
			}
			st.State = DHNeg
		} else {
			tvpn.Sig.SendMessage(Message{Type: Deny, Data: map[string]string{"reason": "Not Authorized"}})
		}
	default:
		tvpn.Sig.SendMessage(Message{Type: Reset, Data: map[string]string{
			"reason": "Invalid state: None"}})
	}
}

// Init state is after Init is sent and before Accept is received
// Next state is DHNeg
func (st *ConState) initState(mes Message, tvpn TVPN) {
	switch mes.Type {
	case Accept:
		st.Params = make([]ovpn.Params, 4)
		st.Key = make([]*big.Int, 4)
		for i := 0; i < 4; i++ {
			st.Params[i] = ovpn.GenParams()
			tvpn.Sig.SendMessage(Message{Type: Dhpub, Data: map[string]string{
				"i": fmt.Sprintf("%d", i),
				"x": base64.StdEncoding.EncodeToString(st.Params[i].X.Bytes()),
				"y": base64.StdEncoding.EncodeToString(st.Params[i].Y.Bytes()),
			}})
		}
		st.State = DHNeg
	default:
		tvpn.Sig.SendMessage(Message{Type: Reset, Data: map[string]string{
			"reason": "Invalid state: Init"}})
	}
}

func (st *ConState) dhnegState(mes Message, tvpn TVPN) {
	switch mes.Type {
	case Dhpub:
		x, y, i, err := mes.DhParams()
		if err != nil {
			tvpn.Sig.SendMessage(Message{Type: Reset, Data: map[string]string{
				"reason": "Invalid DH Params",
			}})
		}
		st.Key[i] = ovpn.GenMutSecret(st.Params[i], ovpn.Params{X: x, Y: y})
		for _, v := range st.Key {
			if v == nil {
				// end state change - still need more keys
				return
			}
		}

	}
}

func (st *ConState) tunnegState(mes Message, tvpn TVPN) {
}

func (st *ConState) connectedState(mes Message, tvpn TVPN) {
}
