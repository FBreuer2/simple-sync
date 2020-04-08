package net_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/FBreuer2/simple-sync/lib/net"
	"github.com/FBreuer2/simple-sync/lib/sync"
)

var helloCombinations = []struct {
	Version        uint16
	Capabilities   uint16
  }{
	{net.VERSION_0_1, net.CAPABILITY_LOGIN | net.CAPABILITY_SYNC | net.CAPABILITY_TOKEN},
  }


func TestHelloMarshalling(t *testing.T) {
	for _, instance := range helloCombinations {
		helloPacket := net.NewHelloPacket()

		marshalled, _ := helloPacket.MarshalBinary()
		helloPacket.UnmarshalBinary(marshalled)

		if helloPacket.Version != instance.Version {
		  t.Errorf("Unmarshaling HelloPacket::Version expected %d, actual %d", instance.Version, helloPacket.Version)
		}

		if helloPacket.Capabilities != instance.Capabilities {
			t.Errorf("Unmarshaling HelloPacket::Capabilities expected %d, actual %d", instance.Capabilities, helloPacket.Capabilities)
		  }
	  }
}


var loginCombinations = []struct {
	username      []byte
	password   []byte
  }{
	{[]byte("admin"), []byte("123")},
	{[]byte("user"), []byte("user")},
  }

func TestLoginPacketMarshalling(t *testing.T) {
	for _, instance := range loginCombinations {
		loginPacket := net.NewLoginPacket(instance.username, instance.password)

		marshalled, _ := loginPacket.MarshalBinary()

		newPacket, _ := net.NewEncapsulatedPacket(loginPacket)

		marshalledPacket, _ := newPacket.MarshalBinary()

		newPacket.UnmarshalBinary(marshalledPacket)
		loginPacket.UnmarshalBinary(newPacket.Data)

		if newPacket.PacketLength != uint64(len(marshalled)) {
			t.Errorf("Unmarshaling packet encapsulated Packet::PacketLength expected %d, actual %d", len(marshalled), newPacket.PacketLength)
		}

		if bytes.Equal(loginPacket.Username, instance.username) != true {
		  t.Errorf("Unmarshaling packet encapsulated LoginPacket::user expected %s, actual %s", string(instance.username), string(loginPacket.Username))
		}

		if bytes.Equal(loginPacket.Password, instance.password) != true {
			t.Errorf("Unmarshaling packet encapsulated LoginPacket::password expected %s, actual %s", string(instance.password), string(loginPacket.Password))
		  }
	  }
}

func TestPacketMarshalling(t *testing.T) {
	for _, instance := range helloCombinations {
		helloPacket := net.NewHelloPacket()

		marshalled, _ := helloPacket.MarshalBinary()

		newPacket, _ := net.NewEncapsulatedPacket(helloPacket)

		marshalledPacket, _ := newPacket.MarshalBinary()

		newPacket.UnmarshalBinary(marshalledPacket)
		helloPacket.UnmarshalBinary(newPacket.Data)

		if newPacket.PacketLength != uint64(len(marshalled)) {
			t.Errorf("Unmarshaling packet encapsulated Packet::PacketLength expected %d, actual %d", len(marshalled), newPacket.PacketLength)
		}

		if helloPacket.Version != instance.Version {
		  t.Errorf("Unmarshaling packet encapsulated HelloPacket::Version expected %d, actual %d", instance.Version, helloPacket.Version)
		}

		if helloPacket.Capabilities != instance.Capabilities {
			t.Errorf("Unmarshaling packet encapsulated HelloPacket::Capabilities expected %d, actual %d", instance.Capabilities, helloPacket.Capabilities)
		  }
	  }
}

var shortFileMetadataCombinations = []*sync.ShortFileMetadata {
	&sync.ShortFileMetadata{12, []byte("123"), time.Now()},
	&sync.ShortFileMetadata{20, []byte("user"), time.Now()},
}

func TestShortFileMetadataPacketMarshalling(t *testing.T) {
	for _, instance := range shortFileMetadataCombinations {
		metaPacket := net.NewShortFileMetaDataPacket(instance)

		marshalled, _ := metaPacket.MarshalBinary()

		newPacket, _ := net.NewEncapsulatedPacket(metaPacket)

		marshalledPacket, _ := newPacket.MarshalBinary()

		newPacket.UnmarshalBinary(marshalledPacket)
		metaPacket.UnmarshalBinary(newPacket.Data)

		if newPacket.PacketLength != uint64(len(marshalled)) {
			t.Errorf("Unmarshaling packet encapsulated Packet::PacketLength expected %d, actual %d", len(marshalled), newPacket.PacketLength)
		}

		if metaPacket.FileSize != instance.FileSize {
			t.Errorf("Unmarshaling packet encapsulated ShortFileMetadataPaket::FileSize expected %d, actual %d", instance.FileSize, metaPacket.FileSize)
		  }

		if bytes.Equal(metaPacket.FileHash, instance.FileHash) != true {
		  t.Errorf("Unmarshaling packet encapsulated ShortFileMetadataPaket::FileHash expected %s, actual %s", string(instance.FileHash), string(metaPacket.FileHash))
		}

		if bytes.Equal(metaPacket.LastChanged, []byte(instance.LastChanged.Format("2006-01-02 15:04:05.999999999 -0700 MST"))) != true {
			t.Errorf("Unmarshaling packet encapsulated ShortFileMetadataPaket::LastChanged expected %s, actual %s", string(instance.LastChanged.Format("2006-01-02 15:04:05.999999999 -0700 MST")), string(metaPacket.LastChanged))
		  }

		_, err := metaPacket.GetData()
		if (err != nil) {
			t.Errorf("Time parsing was errornous: %s", err.Error())
		}
	  }
}