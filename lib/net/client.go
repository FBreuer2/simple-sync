package net

import (
	"net"
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
	client.sendHello()

	for {
		select {
		case <- client.shouldStop:
			client.conn.Close()
			return
		}

	}
}


func (client *ClientContext) sendHello() {
	helloPacket := NewHelloPacket()
	client.sendPacket(helloPacket)

	loginPacket := NewLoginPacket([]byte("user"), []byte("password"))
	client.sendPacket(loginPacket)
}


func (client *ClientContext) sendPacket(packetToSend EncapsulatablePacket) {
	newPacket, _ := NewEncapsulatedPacket(packetToSend)
	data, _ := newPacket.MarshalBinary()
	client.conn.Write(data)
}
