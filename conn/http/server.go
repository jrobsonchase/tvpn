package http

import (
	//"encoding/json"
	"net/http"
	"log"
)

func initHandlers(mux *http.ServeMux,chans TvpnChans) {
	mux.Handle("/init", HandleNew{chans})
	mux.Handle("/term", HandleTerm{chans})
	mux.Handle("/poll", HandlePoll{chans})
	mux.Handle("/connect", HandleConn{chans})
	mux.Handle("/disconnect", HandleDC{chans})
}
func InitServer(addr string) *TvpnServer {
	chans := TvpnChans{
		make(chan map[string] Client,1),
		make(chan map[string] Req,1)}

	log.Print("Putting empty maps...")
	chans.Clients <- map[string] Client {}
	chans.Requests <- map[string] Req {}
	log.Print("Done initializing channels!")

	mux := http.NewServeMux()
	initHandlers(mux,chans)

	log.Print("returning new server")
	return &TvpnServer{chans,&http.Server{
		Addr:    addr,
		Handler: mux}}
}
