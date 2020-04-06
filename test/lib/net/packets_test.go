package net_test

import (
	"bytes"
	"testing"

	"github.com/FBreuer2/simple-sync/lib/net"
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