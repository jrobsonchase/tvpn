package http

import (
	"net/http"
	"encoding/json"
	"bytes"
)

func SendRequest(addr string,req Req) (r *http.Response,err error) {
	data,err := json.Marshal(req)
	if err != nil {
		return
	}
	r,err = http.Post("http://" + addr + "/" + req.Type,"text/json",bytes.NewReader(data))
	return
}
