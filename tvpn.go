package tvpn

import ()

type TVPN struct {
	Name, Group string
	Friends     []string
	Sig         SigBackend
	Stun        StunBackend
	States      map[string]*ConState
}

func New(name, group string,
	friends []string,
	sig SigBackend,
	stun StunBackend,
	vpn VPNBackend) *TVPN {

	return nil
}
