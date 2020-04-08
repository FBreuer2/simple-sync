package net

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"

	"log"
	"net"

	"github.com/FBreuer2/simple-sync/lib/sync"
	"golang.org/x/crypto/sha3"
)

type ClientContext struct {
	url           string
	shouldStop    chan bool
	conn          net.Conn
	authenticated bool
	fileWatcher   *sync.FileWatcher
	serverHash    string
}

func NewClient(url string, serverCertificateHash string) *ClientContext {
	return &ClientContext{
		url:        url,
		serverHash: serverCertificateHash,
	}
}

func (client *ClientContext) checkFingerprint(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	if len(client.serverHash) == 0 {
		return errors.New("No server hash found.")
	}

	decodedServerHash, err := hex.DecodeString(client.serverHash)

	if err != nil {
		return err
	}

	for _, rawCert := range rawCerts {
		if certificateHash := sha3.Sum256(rawCert); bytes.Equal(certificateHash[:], decodedServerHash) == true {
			return nil
		}
	}

	return errors.New("No matching server fingerprint.")
}

func (client *ClientContext) Start(filePath string) error {
	fileWatcher, err := sync.NewFileWatcher(filePath)
	if err != nil {
		log.Fatal(err)
		return err
	}

	client.fileWatcher = fileWatcher

	conf := &tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: client.checkFingerprint,
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
