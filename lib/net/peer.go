package net

import (
	"io"
	"net"
	"log"
)

type Peer struct {
	conn net.Conn
	version uint16
	capabilities uint16
	authenticated bool
	shouldStop chan bool
	closed chan string
}

func NewPeer(conn net.Conn, closed chan string) (*Peer) {
	return &Peer{
		conn: conn,
		shouldStop: make(chan bool),
		closed: closed,
	}
}


func (peer *Peer) GetUniqueIdentifier() (string) {
	return peer.conn.RemoteAddr().String()
}

func (peer *Peer) Start() {
	go peer.mainLoop()

	select {
	case _ = <- peer.shouldStop:
		peer.conn.Close()
		return
	}
}

func (peer *Peer) Stop() {
	peer.shouldStop <- true
}

func (peer *Peer) mainLoop() {
	peekBuf := make([]byte, 10)
	currentAmountOfBytes := 0

	for {
		readNow, err := peer.conn.Read(peekBuf[currentAmountOfBytes:])
		if err != nil {
			if err == io.EOF {
				peer.closed <- peer.GetUniqueIdentifier()
				return
			}

			log.Println(err)
			peer.closed <- peer.GetUniqueIdentifier()
			return
		}

		currentAmountOfBytes += readNow

		// Not enough information to parse a packet yet
		if currentAmountOfBytes < 10 {
			continue
		} else {
			// enough to start parsing
			currentAmountOfBytes = 0
			newPacket, err := PacketFromHeader(peekBuf)

			if err != nil {
				if err == io.EOF {
					peer.closed <- peer.GetUniqueIdentifier()
					return
				}

				log.Println(err)
				peer.closed <- peer.GetUniqueIdentifier()
				return
			}

			packetBuf := make([]byte, newPacket.PacketLength)
			currentPacketAmountOfBytes := uint64(0)

			for {

				readNowPacket, err := peer.conn.Read(packetBuf[currentPacketAmountOfBytes:])

				if err != nil {
					if err == io.EOF {
						peer.closed <- peer.GetUniqueIdentifier()
						return
					}

					log.Println(err)
					peer.closed <- peer.GetUniqueIdentifier()
					return
				}

				currentPacketAmountOfBytes += uint64(readNowPacket)

				if currentPacketAmountOfBytes < newPacket.PacketLength {
					continue
				} else {
					currentPacketAmountOfBytes = 0

					// we can decide which packet it is
					switch newPacket.PacketType {
					case HELLO:
						helloPacket := HelloPacket{}
						helloPacket.UnmarshalBinary(packetBuf)
						peer.HandleHelloPacket(&helloPacket)
						break


					
					}
				}
			}
		}

	}

}


func (peer *Peer) HandleHelloPacket(helloPacket *HelloPacket) {
	peer.version = helloPacket.Version
	peer.capabilities = helloPacket.Capabilities

	log.Printf("Peer on " + peer.conn.RemoteAddr().String() + " sent hello with capability: %d \n", peer.capabilities)
}