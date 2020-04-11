package net

import (
	"encoding/binary"
	"errors"
)

func PacketFromHeader(data []byte) (*Packet, error) {
	if len(data) < 10 {
		return nil, errors.New("PacketFromHeader: Data not long enough")
	}

	return &Packet{
		PacketType:   binary.BigEndian.Uint16(data[:2]),
		PacketLength: binary.BigEndian.Uint64(data[2:10]),
	}, nil
}

func (packet *Packet) MarshalBinary() (data []byte, err error) {
	marshalledData := make([]byte, 10+len(packet.Data))

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

func (rP *ReplyPacket) MarshalBinary() (data []byte, err error) {
	marshalledData := make([]byte, 10+rP.ErrorStringLength)

	binary.BigEndian.PutUint16(marshalledData[:2], rP.ErrorCode)
	binary.BigEndian.PutUint64(marshalledData[2:10], rP.ErrorStringLength)

	copy(marshalledData[10:10+rP.ErrorStringLength], rP.ErrorString)

	return marshalledData, nil
}

func (rP *ReplyPacket) UnmarshalBinary(data []byte) error {
	rP.ErrorCode = binary.BigEndian.Uint16(data[:2])
	rP.ErrorStringLength = binary.BigEndian.Uint64(data[2:10])

	rP.ErrorString = make([]byte, rP.ErrorStringLength)
	copy(rP.ErrorString, data[10:10+rP.ErrorStringLength])
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
	marshalledData := make([]byte, 4+loginPacket.UsernameLength+loginPacket.PasswordLength)

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
	marshalledData := make([]byte, 8+8+8+len(sFM.FileHash)+len(sFM.LastChanged))

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

	sFM.LastChangedLength = binary.BigEndian.Uint64(data[16+sFM.FileHashLength : 24+sFM.FileHashLength])

	sFM.LastChanged = make([]byte, sFM.LastChangedLength)
	copy(sFM.LastChanged, data[24+sFM.FileHashLength:])

	return nil
}

func (eFM *ExtendedFileMetadataPacket) MarshalBinary() (data []byte, err error) {
	marshalledData := make([]byte, 8+4+4+8+eFM.BlockAmount*(4+8)+eFM.BlockAmount*uint64(eFM.StrongChecksumLength))

	binary.BigEndian.PutUint64(marshalledData[:8], eFM.FileSize)
	binary.BigEndian.PutUint32(marshalledData[8:12], eFM.StrongChecksumLength)
	binary.BigEndian.PutUint32(marshalledData[12:16], eFM.BlockLength)
	binary.BigEndian.PutUint64(marshalledData[16:24], eFM.BlockAmount)

	offset := 0
	for key, value := range eFM.WeakBlockHashes {
		binary.BigEndian.PutUint32(marshalledData[24+(offset*12):28+(offset*12)], key)
		binary.BigEndian.PutUint64(marshalledData[28+(offset*12):36+(offset*12)], uint64(value))
		offset += 1
	}

	currentOffset := uint32(36 + ((offset - 1) * 12))
	for index := range eFM.StrongBlockHashes {
		copy(marshalledData[currentOffset:currentOffset+eFM.StrongChecksumLength], eFM.StrongBlockHashes[index])
		currentOffset += eFM.StrongChecksumLength
	}

	return marshalledData, nil
}

func (eFM *ExtendedFileMetadataPacket) UnmarshalBinary(data []byte) error {
	eFM.FileSize = binary.BigEndian.Uint64(data[:8])
	eFM.StrongChecksumLength = binary.BigEndian.Uint32(data[8:12])
	eFM.BlockLength = binary.BigEndian.Uint32(data[12:16])
	eFM.BlockAmount = binary.BigEndian.Uint64(data[16:24])

	eFM.WeakBlockHashes = make(map[uint32]int64)

	offset := uint64(0)

	for offset < eFM.BlockAmount {
		key := binary.BigEndian.Uint32(data[24+(offset*12) : 28+(offset*12)])
		value := binary.BigEndian.Uint64(data[28+(offset*12) : 36+(offset*12)])

		eFM.WeakBlockHashes[key] = int64(value)
		offset += 1

	}

	eFM.StrongBlockHashes = make([][]byte, eFM.BlockAmount)

	currentOffset := uint32(36 + ((offset - 1) * 12))

	for index := range eFM.StrongBlockHashes {
		eFM.StrongBlockHashes[index] = make([]byte, eFM.StrongChecksumLength)
		copy(eFM.StrongBlockHashes[index], data[currentOffset:currentOffset+eFM.StrongChecksumLength])
		currentOffset += eFM.StrongChecksumLength
	}

	return nil
}
