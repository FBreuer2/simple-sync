package net

import (
	"net"
    "log"
	"crypto/tls"
)

type ClientContext struct {
	url string
	shouldStop chan bool
	conn net.Conn
}

func NewClient(url string) (*ClientContext) {
	return &ClientContext{
		url: url,
	}

	
}

func (client *ClientContext) Start() (error) {
   	conf := &tls.Config{
		InsecureSkipVerify: true,
   	}

   	newConnection, err := tls.Dial("tcp", client.url, conf)
   	if err != nil {
	   	log.Println(err)
	   	return err
	}

	client.conn = newConnection
	   
	go client.mainLoop()

	return nil
}

func (client *ClientContext) Stop() {
	client.shouldStop <- true
}


func (client *ClientContext) mainLoop() {
	helloPacket := NewHelloPacket()

	newPacket, _ := NewEncapsulatedPacket(helloPacket)

	data, _ := newPacket.MarshalBinary()
	client.conn.Write(data)
}
