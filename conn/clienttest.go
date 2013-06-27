package main

import (
	"tvpn/conn/http"
	//"fmt"
)

func main() {
	req := http.Req{To: "localhost",Type: http.Init,From: http.Client{Name: "test2",IP: "127.0.0.1",Port: "12345",Key: "KEY!"}}
	http.SendRequest("localhost:1234",req)
}
