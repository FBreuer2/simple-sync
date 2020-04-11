package net

import (
	"encoding"
	"time"

	"github.com/FBreuer2/simple-sync/lib/sync"
)

const (
	REPLY                  = 0
	HELLO                  = 1
	LOGIN                  = 2
	SHORT_FILE_METADATA    = 3
	EXTENDED_FILE_METADATA = 4
)

const (
	VERSION_0_1 = 0
)

const (
	CAPABILITY_LOGIN = 0
	CAPABILITY_SYNC  = 1
	CAPABILITY_TOKEN = 2
)

type Packet struct {
	PacketType   uint16
	PacketLength uint64
	Data         []byte
}

func NewEncapsulatedPacket(originalPacket EncapsulatablePacket) (*Packet, error) {
	marshalled, err := originalPacket.MarshalBinary()

	if err != nil {
		return nil, err
	}

	return &Packet{originalPacket.Type(), uint64(len(marshalled)), marshalled}, nil
}

type EncapsulatablePacket interface {
	Type() uint16
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type ReplyPacket struct {
	ErrorCode         uint16
	ErrorStringLength uint64
	ErrorString       []byte
}

func NewReplyPacket(code uint16, errorString string) *ReplyPacket {
	return &ReplyPacket{code, uint64(len([]byte(errorString))), []byte(errorString)}
}

func (rp *ReplyPacket) Type() uint16 {
	return REPLY
}

type HelloPacket struct {
	Version      uint16
	Capabilities uint16
}

func NewHelloPacket() *HelloPacket {
	return &HelloPacket{VERSION_0_1, CAPABILITY_LOGIN | CAPABILITY_SYNC | CAPABILITY_TOKEN}
}

func (hP *HelloPacket) Type() uint16 {
	return HELLO
}

type LoginPacket struct {
	UsernameLength uint16
	Username       []byte
	PasswordLength uint16
	Password       []byte
}

func NewLoginPacket(username, password []byte) *LoginPacket {
	return &LoginPacket{
		UsernameLength: uint16(len(username)),
		Username:       username,
		PasswordLength: uint16(len(password)),
		Password:       password,
	}
}

func (lP *LoginPacket) Type() uint16 {
	return LOGIN
}

type ShortFileMetadataPacket struct {
	FileSize          uint64
	FileHashLength    uint64
	FileHash          []byte
	LastChangedLength uint64
	LastChanged       []byte // String() string
}

func NewShortFileMetaDataPacket(sFM *sync.ShortFileMetadata) *ShortFileMetadataPacket {
	return &ShortFileMetadataPacket{
		FileSize:          sFM.FileSize,
		FileHashLength:    uint64(len(sFM.FileHash)),
		FileHash:          sFM.FileHash,
		LastChangedLength: uint64(len([]byte(sFM.LastChanged.Format("2006-01-02 15:04:05.999999999 -0700 MST")))),
		LastChanged:       []byte(sFM.LastChanged.Format("2006-01-02 15:04:05.999999999 -0700 MST")),
	}
}

func (sFM *ShortFileMetadataPacket) GetData() (*sync.ShortFileMetadata, error) {
	timeObject, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(sFM.LastChanged))

	if err != nil {
		return nil, err
	}

	return &sync.ShortFileMetadata{
		FileSize:    sFM.FileSize,
		FileHash:    sFM.FileHash,
		LastChanged: timeObject,
	}, nil
}

func (sFM *ShortFileMetadataPacket) Type() uint16 {
	return SHORT_FILE_METADATA
}

type ExtendedFileMetadataPacket struct {
	FileSize             uint64
	StrongChecksumLength uint32
	BlockLength          uint32
	BlockAmount          uint64
	WeakBlockHashes      map[uint32]int64
	StrongBlockHashes    [][]byte
}

func NewExtendedFileMetadataPacket(eFM *sync.ExtendedFileMetadata) (*ExtendedFileMetadataPacket, error) {
	return &ExtendedFileMetadataPacket{
		FileSize:             eFM.FileSize,
		StrongChecksumLength: eFM.StrongChecksumLength,
		BlockLength:          eFM.BlockLength,
		BlockAmount:          eFM.BlockAmount,
		WeakBlockHashes:      eFM.WeakBlockHashes,
		StrongBlockHashes:    eFM.StrongBlockHashes,
	}, nil
}

func (eFM *ExtendedFileMetadataPacket) GetData() (*sync.ExtendedFileMetadata, error) {
	return &sync.ExtendedFileMetadata{
		FileSize:             eFM.FileSize,
		StrongChecksumLength: eFM.StrongChecksumLength,
		BlockLength:          eFM.BlockLength,
		BlockAmount:          eFM.BlockAmount,
		WeakBlockHashes:      eFM.WeakBlockHashes,
		StrongBlockHashes:    eFM.StrongBlockHashes,
	}, nil
}

func (eFM *ExtendedFileMetadataPacket) Type() uint16 {
	return EXTENDED_FILE_METADATA
}
