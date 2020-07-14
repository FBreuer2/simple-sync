package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/FBreuer2/simple-sync/lib/net"
)

func main() {

	var pathToInputFile, serverURL, fingerprint string

	flag.StringVar(&serverURL, "u", "127.0.0.1:8888", "URL of the server.")
	flag.StringVar(&pathToInputFile, "i", "./file.bmp", "Path to a file which should be version controlled.")
	flag.StringVar(&fingerprint, "f", "2926ea1c1e4adefb2ecbc7bb58e3172752d36274fdf899aca1667debd292fb7b", "Fingerprint of the server's tls certificate.")
	flag.Parse()

	client := net.NewClient(serverURL, fingerprint)
	err := client.Start(pathToInputFile)

	if err != nil {
		log.Println(err)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		go client.Stop()
		return
	}
}
