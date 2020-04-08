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


func (loginPacket *LoginPacket) MarshalBinary() (data []byte, err error) {
	marshalledData := make([]byte, 4 + loginPacket.UsernameLength + loginPacket.PasswordLength)

	binary.BigEndian.PutUint16(marshalledData[:2], loginPacket.UsernameLength)
	copy(marshalledData[2:2+loginPacket.UsernameLength], loginPacket.Username)

	binary.BigEndian.PutUint16(marshalledData[2+loginPacket.UsernameLength:], loginPacket.PasswordLength)
	copy(marshalledData[4+loginPacket.UsernameLength:], loginPacket.Password)

	return marshalledData, nil
}


func (loginPacket *LoginPacket) UnmarshalBinary(data []byte) error {
	loginPacket.UsernameLength = binary.BigEndian.Uint16(data[:2])

	loginPacket.Username = make([]byte, loginPacket.UsernameLength)
	copy(loginPacket.Username, data[2:2+loginPacket.UsernameLength])

	loginPacket.PasswordLength = binary.BigEndian.Uint16(data[2+loginPacket.UsernameLength:])

	loginPacket.Password = make([]byte, loginPacket.PasswordLength)
	copy(loginPacket.Password, data[4+loginPacket.UsernameLength:])

	return nil
}


func (sFM *ShortFileMetadataPacket) MarshalBinary() (data []byte, err error) {
	marshalledData := make([]byte, 8 + 8 + 8  + len(sFM.FileHash) + len(sFM.LastChanged))

	binary.BigEndian.PutUint64(marshalledData[:8], sFM.FileSize)
	binary.BigEndian.PutUint64(marshalledData[8:16], sFM.FileHashLength)

	copy(marshalledData[16:16+sFM.FileHashLength], sFM.FileHash)

	binary.BigEndian.PutUint64(marshalledData[16+sFM.FileHashLength:24+sFM.FileHashLength], sFM.LastChangedLength)
	copy(marshalledData[24+sFM.FileHashLength:24+sFM.FileHashLength+sFM.LastChangedLength], sFM.LastChanged)

	return marshalledData, nil
}


func (sFM *ShortFileMetadataPacket) UnmarshalBinary(data []byte) error {
	sFM.FileSize = binary.BigEndian.Uint64(data[:8])

	sFM.FileHashLength = binary.BigEndian.Uint64(data[8:16])

	sFM.FileHash = make([]byte, sFM.FileHashLength)
	copy(sFM.FileHash, data[16:16+sFM.FileHashLength])

	sFM.LastChangedLength = binary.BigEndian.Uint64(data[16+sFM.FileHashLength:24+sFM.FileHashLength])

	sFM.LastChanged = make([]byte, sFM.LastChangedLength)
	copy(sFM.LastChanged, data[24+sFM.FileHashLength:])

	return nil
}