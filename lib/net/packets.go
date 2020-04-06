package net

import (
	"encoding"
)

const (
	HELLO = 1
	LOGIN = 2
)

const (
	VERSION_0_1    = 0
)

const (
	CAPABILITY_LOGIN = 0
	CAPABILITY_SYNC = 1
	CAPABILITY_TOKEN = 2
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
	return &HelloPacket{VERSION_0_1, CAPABILITY_LOGIN | CAPABILITY_SYNC | CAPABILITY_TOKEN}
}


func ( hP*HelloPacket) Type() uint16 {
	return HELLO
}


type LoginPacket struct {
	UsernameLength uint16
	Username       []byte
	PasswordLength uint16
	Password       []byte
}


func NewLoginPacket(username, password []byte) (*LoginPacket) {
	return &LoginPacket{
		UsernameLength: uint16(len(username)),
		Username: username,
		PasswordLength: uint16(len(password)),
		Password: password,
	}
}

func (lP *LoginPacket) Type() uint16 {
	return LOGIN
}
