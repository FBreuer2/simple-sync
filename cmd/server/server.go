package main

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"github.com/FBreuer2/simple-sync/lib/net"
)

func main() {

	cer, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
    if err != nil {
        log.Println(err)
        return
	}

	srv, err := net.NewServer("127.0.0.1", "8888", cer)

    if err != nil {
        log.Println(err)
        return
	}

	err = srv.Start()

	if err != nil {
        log.Println(err)
        return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
		case <-c:
			srv.Stop()
		return
	}
}