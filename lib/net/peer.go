package net

import (
	"io"
	"log"
	"net"

	"github.com/FBreuer2/simple-sync/lib/db"
)

type Peer struct {
	conn          net.Conn
	version       uint16
	capabilities  uint16
	authenticated bool
	username      []byte
	shouldStop    chan bool
	closed        chan string
	db            db.FullDatabase
}

func NewPeer(conn net.Conn, closed chan string, db db.FullDatabase) *Peer {
	return &Peer{
		conn:       conn,
		shouldStop: make(chan bool),
		closed:     closed,
		db:         db,
	}
}

func (peer *Peer) GetUniqueIdentifier() string {
	return peer.conn.RemoteAddr().String()
}

func (peer *Peer) Start() {
	go peer.mainLoop()

	select {
	case _ = <-peer.shouldStop:
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

					case LOGIN:
						loginPacket := LoginPacket{}
						loginPacket.UnmarshalBinary(packetBuf)
						peer.HandleLoginPacket(&loginPacket)
						break

					case SHORT_FILE_METADATA:
						if peer.authenticated == false {
							// XXX: send error
							log.Printf("Peer on " + peer.conn.RemoteAddr().String() + " sent packet without authentication\n")
							break
						}

						sFMPPacket := ShortFileMetadataPacket{}
						sFMPPacket.UnmarshalBinary(packetBuf)
						peer.HandleShortFileMetadataPacketPacket(&sFMPPacket)
						break
					}

					break
				}
			}
		}

	}

}

func (peer *Peer) HandleHelloPacket(helloPacket *HelloPacket) {
	peer.version = helloPacket.Version
	peer.capabilities = helloPacket.Capabilities

	log.Printf("Peer on "+peer.conn.RemoteAddr().String()+" sent hello with capability: %d \n", peer.capabilities)
}

func (peer *Peer) HandleLoginPacket(loginPacket *LoginPacket) {
	err := peer.db.Login(loginPacket.Username, loginPacket.Password)

	if err != nil {
		// XXX: send error
		log.Printf("Peer on "+peer.conn.RemoteAddr().String()+" tried to authenticate for \"%s\" with error: %s\n", string(loginPacket.Username), err.Error())
		return
	}

	peer.authenticated = true
	peer.username = loginPacket.Username

	log.Printf("Peer on "+peer.conn.RemoteAddr().String()+" authenticated for \"%s\" \n", string(loginPacket.Username))
}

func (peer *Peer) HandleShortFileMetadataPacketPacket(shortFileMetadataPacket *ShortFileMetadataPacket) {
	newSFM, err := shortFileMetadataPacket.GetData()

	if err != nil {
		log.Printf("Peer on "+peer.conn.RemoteAddr().String()+" has error \"%s\" \n", err.Error())
		return
	}

	currentSFM, err := peer.db.RetrieveShortFileMetadata(peer.username)

	if err != nil || newSFM.ShouldOverwrite(currentSFM) == true {
		// SFM not saved yet
		peer.db.PutShortFileMetadata(peer.username, newSFM)
		log.Printf("Peer on "+peer.conn.RemoteAddr().String()+" sent new metadata with file size %d and time %s\n", newSFM.FileSize, newSFM.LastChanged.Format("2006-01-02 15:04:05.999999999 -0700 MST"))
		go peer.RetrieveBlocks()
		return
	}

	// stale metadata
	if newSFM.Equals(currentSFM) == false {
		log.Printf("Peer on "+peer.conn.RemoteAddr().String()+" has stale file with file size %d and time %s\n", newSFM.FileSize, newSFM.LastChanged.Format("2006-01-02 15:04:05.999999999 -0700 MST"))
		return
	}

	log.Printf("Peer on "+peer.conn.RemoteAddr().String()+" sent same metadata with file size %d and time %s\n", newSFM.FileSize, newSFM.LastChanged.Format("2006-01-02 15:04:05.999999999 -0700 MST"))
	return
}

func (peer *Peer) RetrieveBlocks() {
	// Check which blocks we have

	// Send them to the client
}
