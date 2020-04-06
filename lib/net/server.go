package net

import (
	"net"
	"crypto/tls"
	"log"
)

type ServerContext struct {
	interfaceToBind string
	port string
	cert tls.Certificate
	shouldStop chan bool
	closed chan string
	acceptingServer net.Listener
	peerList map[string]*Peer
}


func NewServer(interfaceToBind string, port string, cert tls.Certificate) (*ServerContext, error) {
	newServerContext := &ServerContext{
		interfaceToBind: interfaceToBind,
		port: port,
		cert: cert,
		shouldStop: make(chan bool),
		closed: make(chan string),
		peerList: make(map[string]*Peer),
	}

	return newServerContext, nil
}


func (srv *ServerContext) Start() (error) {
	config := &tls.Config{
		MinVersion:               tls.VersionTLS12,
        CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
        },
		Certificates: []tls.Certificate{srv.cert},
	}

	newListener, err := tls.Listen("tcp", srv.interfaceToBind + ":" + srv.port, config) 

	log.Println("Started server.")

	srv.acceptingServer = newListener

    if err != nil {
        return err
	}

	go srv.mainLoop()

	return nil
}


func (srv *ServerContext) Stop() {
	srv.shouldStop <- true
}

func (srv *ServerContext)  mainLoop() {
	log.Println("Started main loop.")

	go srv.runAccept()
	
	for {
		select {
		case peerID := <- srv.closed:
			srv.peerList[peerID] = nil
			log.Println("Peer with id " + peerID + " disconnected.")
			break
		case <- srv.shouldStop:
			for _, peer := range(srv.peerList) {
				peer.Stop()
			}

			srv.acceptingServer.Close()
			return
		}
	}

}

func (srv *ServerContext) runAccept() {
	for {
		conn, err := srv.acceptingServer.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
	
		srv.newClient(conn)
	}
}


func (srv *ServerContext) newClient(newClient net.Conn) {
	newPeer := NewPeer(newClient, srv.closed)

	exists := srv.peerList[newPeer.GetUniqueIdentifier()]
	
	if (exists != nil) {
		log.Println("Unique ID " + newPeer.GetUniqueIdentifier() + " is already in use.")
		return
	}

	srv.peerList[newPeer.GetUniqueIdentifier()] = newPeer
	go newPeer.Start()
}
