package net

import (
	"crypto/tls"

	"log"
	"net"

	"github.com/FBreuer2/simple-sync/lib/sync"
)

type ClientContext struct {
	url         string
	shouldStop  chan bool
	conn        net.Conn
	fileWatcher *sync.FileWatcher
}

func NewClient(url string) *ClientContext {
	return &ClientContext{
		url: url,
	}
}

func (client *ClientContext) Start(filePath string) error {
	fileWatcher, err := sync.NewFileWatcher(filePath)
	if err != nil {
		log.Fatal(err)
		return err
	}

	client.fileWatcher = fileWatcher

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
	client.sendLoginPacket()
	client.sendShortFileMetadata()

	for {
		select {
		case <-client.shouldStop:
			client.conn.Close()
			return
		}

	}
}

func (client *ClientContext) sendHello() {
	helloPacket := NewHelloPacket()
	err := client.sendPacket(helloPacket)

	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (client *ClientContext) sendShortFileMetadata() {
	shortFileMetadata, err := client.fileWatcher.GetShortFileMetadata()

	if err != nil {
		log.Println(err.Error())
		return
	}

	shortFileMetadataPacket := NewShortFileMetaDataPacket(shortFileMetadata)

	err = client.sendPacket(shortFileMetadataPacket)

	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}

func (client *ClientContext) sendLoginPacket() {
	loginPacket := NewLoginPacket([]byte("user"), []byte("password"))

	err := client.sendPacket(loginPacket)

	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (client *ClientContext) sendPacket(packetToSend EncapsulatablePacket) error {
	newPacket, err := NewEncapsulatedPacket(packetToSend)
	if err != nil {
		return err
	}

	data, err := newPacket.MarshalBinary()
	if err != nil {
		return err
	}

	client.conn.Write(data)

	return nil
}
