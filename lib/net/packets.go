package net

import (
	"encoding"
)

const (
	HELLO = 0
)

const (
	VERSION_0    = 0

	CAPABILITY_DELTA = 0
)


type Packet struct {
	PacketType   uint16
	PacketLength uint64
	Data         []byte
}

func NewEncapsulatedPacket(originalPacket EncapsulatablePacket) (*Packet, error) {
	marshalled, err := originalPacket.MarshalBinary()

	if (err != nil) {
		return nil, err
	}

	return &Packet{originalPacket.Type(), uint64(len(marshalled)), marshalled}, nil
}

type EncapsulatablePacket interface {
	Type() 				uint16
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type HelloPacket struct {
	Version      uint16
	Capabilities uint16
}

func NewHelloPacket() (*HelloPacket) {
	return &HelloPacket{VERSION_0, CAPABILITY_DELTA}
}

func ( hP*HelloPacket) Type() uint16 {
	return HELLO
}