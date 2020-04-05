package net

import (
	"encoding/binary"
	"errors"
)

func PacketFromHeader(data []byte) (*Packet, error) { 
	if (len(data) < 10) {
		return nil, errors.New("PacketFromHeader: Data not long enough")
	}

	return &Packet{
		PacketType: binary.BigEndian.Uint16(data[:2]),
		PacketLength: binary.BigEndian.Uint64(data[2:10]),
	}, nil
}

func (packet *Packet) MarshalBinary() (data []byte, err error) {
	marshalledData := make([]byte, 10 + len(packet.Data))

	binary.BigEndian.PutUint16(marshalledData[:2], packet.PacketType)
	binary.BigEndian.PutUint64(marshalledData[2:10], packet.PacketLength)
	copy(marshalledData[10:], packet.Data) 

	return marshalledData, nil
}


func (packet *Packet) UnmarshalBinary(data []byte) error {
	packet.PacketType = binary.BigEndian.Uint16(data[:2])
	packet.PacketLength = binary.BigEndian.Uint64(data[2:10])

	copy(packet.Data, data[10:])
	
	return nil
}

func (helloPacket *HelloPacket) MarshalBinary() (data []byte, err error) {
	marshalledData := make([]byte, 4)

	binary.BigEndian.PutUint16(marshalledData[:2], helloPacket.Version)
	binary.BigEndian.PutUint16(marshalledData[2:], helloPacket.Capabilities)

	return marshalledData, nil
}


func (helloPacket *HelloPacket) UnmarshalBinary(data []byte) error {
	helloPacket.Version = binary.BigEndian.Uint16(data[:2])
	helloPacket.Capabilities = binary.BigEndian.Uint16(data[2:])

	return nil
}