package http

import (
	"log"
	"io"
	"encoding/json"
	"net/http"
)

type BufferFull bool
func (b *BufferFull) Error() string {
	return "Buffer filled when reading!"
}

func trimNull(bs []byte) []byte {
	for i,v := range bs {
		if v == '\x00' {
			return bs[:i]
		}
	}
	return bs
}

func readRequest(r io.Reader) (*Req,error) {
	buf := make([]byte,1024)
	_,err := io.ReadFull(r,buf)
	if err != io.ErrUnexpectedEOF {
		if err != nil {
			return nil,err
		} else {
			errFull := BufferFull(true)
			return nil,&errFull
		}
	}
	buf = trimNull(buf)
	var ret *Req = new(Req)
	err = json.Unmarshal(buf,ret)
	return ret,err
}

func PrintReq(req Req) {
	PrintClient(req.From)
	log.Println(req.To)
	log.Println(req.Type)
}
func PrintClient(cl Client) {
	log.Println(cl.Name)
	log.Println(cl.IP)
	log.Println(cl.Port)
	log.Println(cl.Key)
}

func (h HandleNew) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	req,err := readRequest(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	PrintReq(*req)
	if req.Type != Init {
		log.Printf("Error: invalid Req type from %s",r.RemoteAddr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	clients := <-h.Clients
	_,exists := clients[req.From.Name]
	if exists {
		log.Printf("Client %s already exists!",req.From.Name)
		w.WriteHeader(http.StatusNotFound)
	} else {
		clients[req.From.Name] = req.From
		w.WriteHeader(http.StatusOK)
	}
	h.Clients <- clients
}

func (h HandleTerm) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	req,err := readRequest(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Type != Term {
		log.Printf("Error: invalid Req type from %s",r.RemoteAddr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	clients := <-h.Clients
	_,exists := clients[req.From.Name]
	if ! exists {
		log.Printf("Client %s does not exist!",req.From.Name)
		w.WriteHeader(http.StatusNotFound)
	} else {
		delete(clients,req.From.Name)
		w.WriteHeader(http.StatusOK)
	}
	h.Clients <- clients
}
func (h HandlePoll) ServeHTTP(w http.ResponseWriter,r *http.Request) {
}
func (h HandleConn) ServeHTTP(w http.ResponseWriter,r *http.Request) {
}
func (h HandleDC) ServeHTTP(w http.ResponseWriter,r *http.Request) {
}
