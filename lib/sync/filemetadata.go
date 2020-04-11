package sync

import (
	"bytes"
	"time"
)

type ExtendedFileMetadata struct {
	FileSize             uint64
	StrongChecksumLength uint32
	BlockLength          uint32
	BlockAmount          uint64
	WeakBlockHashes      map[uint32]int64
	StrongBlockHashes    [][]byte
}

func (eFM *ExtendedFileMetadata) Equals(otherEFM *ExtendedFileMetadata) bool {
	if len(otherEFM.StrongBlockHashes) != len(eFM.StrongBlockHashes) {
		return false
	}

	for index := range eFM.StrongBlockHashes {
		if bytes.Equal(eFM.StrongBlockHashes[index], otherEFM.StrongBlockHashes[index]) == false {
			return false
		}
	}

	if len(otherEFM.WeakBlockHashes) != len(eFM.WeakBlockHashes) {
		return false
	}

	for key, value := range eFM.WeakBlockHashes {
		if otherEFM.WeakBlockHashes[key] != value {
			return false
		}
	}

	if eFM.FileSize != otherEFM.FileSize ||
		eFM.StrongChecksumLength != otherEFM.StrongChecksumLength ||
		eFM.BlockLength != otherEFM.BlockLength ||
		eFM.BlockAmount != otherEFM.BlockAmount {
		return false
	}

	return true
}

type ShortFileMetadata struct {
	FileSize    uint64
	FileHash    []byte
	LastChanged time.Time // String() string
}

func (sFM *ShortFileMetadata) ShouldOverwrite(otherSFM *ShortFileMetadata) bool {
	return (sFM.LastChanged.After(otherSFM.LastChanged) && bytes.Equal(sFM.FileHash, otherSFM.FileHash) == false)
}

func (sFM *ShortFileMetadata) Equals(otherSFM *ShortFileMetadata) bool {
	return (sFM.FileSize == otherSFM.FileSize && bytes.Equal(sFM.FileHash, otherSFM.FileHash))
}
