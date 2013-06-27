package main

import (
	"tvpn/conn/http"
	"fmt"
	"log"
	"time"
)

func main() {
	serv := http.InitServer(":1234")
	fmt.Print("Serving http...\n")
	go func() {
		err := serv.HTTP.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	for {
		time.Sleep(10 * time.Second)
		cls := <-serv.Clients
		for _,v := range cls {
			http.PrintClient(v)
		}
		serv.Clients <-cls
	}
}
