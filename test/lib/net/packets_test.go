package net_test

import (
	"testing"
	"github.com/FBreuer2/simple-sync/lib/net"
)

var combinations = []struct {
	Version        uint16
	Capabilities   uint16
  }{
	{net.VERSION_0, net.CAPABILITY_DELTA},
  }


func TestHelloMarshalling(t *testing.T) {
	for _, instance := range combinations {
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

func TestPacketMarshalling(t *testing.T) {
	for _, instance := range combinations {
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