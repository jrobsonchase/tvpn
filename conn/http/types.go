package http

import "net/http"

const (
	Init string = "init"
	Term string = "term"
	Connect string = "connect"
	Disconnect string = "disconnect"
	Poll string = "poll"
)

type Client struct {
	Name string
	IP string
	Port string
	Key string
}

type Req struct {
	From Client
	To string
	Type string
}

type TvpnServer struct {
	TvpnChans
	HTTP *http.Server
}

type TvpnChans struct {
	Clients chan map[string] Client
	Requests chan map[string] Req
}

type HandleNew struct {
	TvpnChans
}
type HandleTerm struct {
	TvpnChans
}
type HandlePoll struct {
	TvpnChans
}
type HandleConn struct {
	TvpnChans
}
type HandleDC struct {
	TvpnChans
}
