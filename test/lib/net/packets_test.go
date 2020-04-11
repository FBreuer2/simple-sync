package net_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/FBreuer2/simple-sync/lib/net"
	"github.com/FBreuer2/simple-sync/lib/sync"
)

var helloCombinations = []struct {
	Version      uint16
	Capabilities uint16
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
	username []byte
	password []byte
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

var shortFileMetadataCombinations = []*sync.ShortFileMetadata{
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
		if err != nil {
			t.Errorf("Time parsing was errornous: %s", err.Error())
		}
	}
}

var errorCombinations = []struct {
	errorCode         uint16
	errorStringLength uint64
	errorString       []byte
}{
	{1, uint64(len([]byte("abcde"))), []byte("abcde")},
	{2, uint64(len([]byte("def"))), []byte("def")},
	{0, uint64(len([]byte(""))), []byte("")},
}

func TestReplyPacketMarshalling(t *testing.T) {
	for _, instance := range errorCombinations {
		replyPacket := net.NewReplyPacket(instance.errorCode, string(instance.errorString))

		marshalled, _ := replyPacket.MarshalBinary()

		newPacket, _ := net.NewEncapsulatedPacket(replyPacket)

		marshalledPacket, _ := newPacket.MarshalBinary()

		newPacket.UnmarshalBinary(marshalledPacket)
		replyPacket.UnmarshalBinary(newPacket.Data)

		if newPacket.PacketLength != uint64(len(marshalled)) {
			t.Errorf("Unmarshaling packet encapsulated Packet::PacketLength expected %d, actual %d", len(marshalled), newPacket.PacketLength)
		}

		if replyPacket.ErrorCode != instance.errorCode {
			t.Errorf("Unmarshaling packet encapsulated ReplyPacket::ErrorCode expected %d, actual %d", instance.errorCode, replyPacket.ErrorCode)
		}

		if replyPacket.ErrorStringLength != instance.errorStringLength {
			t.Errorf("Unmarshaling packet encapsulated ReplyPacket::ErrorStringLength expected %d, actual %d", instance.errorStringLength, replyPacket.ErrorStringLength)
		}

		if bytes.Equal(replyPacket.ErrorString, instance.errorString) != true {
			t.Errorf("Unmarshaling packet encapsulated ShortFileMetadataPaket::FileHash expected %s, actual %s", string(instance.errorString), string(replyPacket.ErrorString))
		}
	}
}

var extendedFileMetadataCombinations = []*sync.ExtendedFileMetadata{
	&sync.ExtendedFileMetadata{12, 2, 3, 5, make(map[uint32]int64), make([][]byte, 5)},
	&sync.ExtendedFileMetadata{342, 23, 3, 5, make(map[uint32]int64), make([][]byte, 5)},
}

func TestExtendedFileMetadataPacketMarshalling(t *testing.T) {
	for _, instance := range extendedFileMetadataCombinations {

		for i := 0; uint64(i) < instance.BlockAmount; i++ {
			instance.WeakBlockHashes[uint32(i)] = int64(i + 1)
		}

		for index := range instance.StrongBlockHashes {
			instance.StrongBlockHashes[index] = make([]byte, instance.StrongChecksumLength)

			for innerIndex := range instance.StrongBlockHashes[index] {
				instance.StrongBlockHashes[index][innerIndex] = 'q'
			}
		}

		metaPacket, _ := net.NewExtendedFileMetadataPacket(instance)

		marshalled, _ := metaPacket.MarshalBinary()

		newPacket, _ := net.NewEncapsulatedPacket(metaPacket)

		marshalledPacket, _ := newPacket.MarshalBinary()

		newPacket.UnmarshalBinary(marshalledPacket)
		metaPacket.UnmarshalBinary(newPacket.Data)

		if newPacket.PacketLength != uint64(len(marshalled)) {
			t.Errorf("Unmarshaling packet encapsulated Packet::PacketLength expected %d, actual %d", len(marshalled), newPacket.PacketLength)
		}

		if metaPacket.FileSize != instance.FileSize {
			t.Errorf("Unmarshaling packet encapsulated ExtendedFileMetadataPacket::FileSize expected %d, actual %d", instance.FileSize, metaPacket.FileSize)
		}

		if metaPacket.StrongChecksumLength != instance.StrongChecksumLength {
			t.Errorf("Unmarshaling packet encapsulated ExtendedFileMetadataPacket::StrongChecksumLength expected %d, actual %d", instance.StrongChecksumLength, metaPacket.StrongChecksumLength)
		}

		if metaPacket.BlockLength != instance.BlockLength {
			t.Errorf("Unmarshaling packet encapsulated ExtendedFileMetadataPacket::BlockLength expected %d, actual %d", instance.BlockLength, metaPacket.BlockLength)
		}

		if metaPacket.BlockAmount != instance.BlockAmount {
			t.Errorf("Unmarshaling packet encapsulated ExtendedFileMetadataPacket::BlockAmount expected %d, actual %d", instance.BlockAmount, metaPacket.BlockAmount)
		}

		retrieved, _ := metaPacket.GetData()
		if retrieved.Equals(instance) == false {
			t.Errorf("ExtendedFileMetadataPacket::Equals failed")
		}
	}
}

var blockPacketCombinations = []*net.BlockPacket{
	&net.BlockPacket{4, []byte("abcd"), 6, []byte("dcefad")},
	&net.BlockPacket{5, []byte("abcda"), 7, []byte("dcesfad")},
}

func TestBlockPacketMarshalling(t *testing.T) {
	for _, instance := range blockPacketCombinations {
		metaPacket := instance

		marshalled, _ := metaPacket.MarshalBinary()

		newPacket, _ := net.NewEncapsulatedPacket(metaPacket)

		marshalledPacket, _ := newPacket.MarshalBinary()

		newPacket.UnmarshalBinary(marshalledPacket)
		metaPacket.UnmarshalBinary(newPacket.Data)

		if newPacket.PacketLength != uint64(len(marshalled)) {
			t.Errorf("Unmarshaling packet encapsulated Packet::PacketLength expected %d, actual %d", len(marshalled), newPacket.PacketLength)
		}

		if metaPacket.Equals(instance) == false {
			t.Errorf("BlockPacket::Equals failed")
		}
	}
}

var requestBlockPacketCombinations = []*net.RequestBlockPacket{
	&net.RequestBlockPacket{4, []byte("abcd")},
	&net.RequestBlockPacket{5, []byte("abcda")},
}

func TestRequestBlockPacketMarshalling(t *testing.T) {
	for _, instance := range requestBlockPacketCombinations {
		metaPacket := instance

		marshalled, _ := metaPacket.MarshalBinary()

		newPacket, _ := net.NewEncapsulatedPacket(metaPacket)

		marshalledPacket, _ := newPacket.MarshalBinary()

		newPacket.UnmarshalBinary(marshalledPacket)
		metaPacket.UnmarshalBinary(newPacket.Data)

		if newPacket.PacketLength != uint64(len(marshalled)) {
			t.Errorf("Unmarshaling packet encapsulated Packet::PacketLength expected %d, actual %d", len(marshalled), newPacket.PacketLength)
		}

		if metaPacket.Equals(instance) == false {
			t.Errorf("RequestBlockPacket::Equals failed")
		}
	}
}
