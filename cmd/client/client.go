package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/FBreuer2/simple-sync/lib/net"
)

var pathToInputFile, serverURL string
var blockLength, strongLength uint

func main() {

	flag.StringVar(&serverURL, "u", "127.0.0.1:8888", "URL of the server.")
	flag.StringVar(&pathToInputFile, "i", "/path/to/file", "Path to a file which should be version controlled.")
	flag.UintVar(&blockLength, "b", 10, "Length of a block in kilobytes")
	flag.UintVar(&strongLength, "c", 16, "Length of the hash value of the strong checksum in bytes")
	flag.Parse()

	client := net.NewClient(serverURL)
	err := client.Start("./keys.kdbx")

	if (err != nil) {
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
