package sync

import (
	"time"
)

type ExtendedFileMetadata struct {
	FileSize             uint64
	StrongChecksumLength uint32
	BlockLength          uint32
	BlockAmount          uint64
	WeakBlockHashes      map[uint32]int
	StrongBlockHashes    [][]byte
}

type ShortFileMetadata struct {
	FileSize    uint64
	FileHash    []byte
	LastChanged time.Time // String() string
}

func (sFM *ShortFileMetadata) ShouldOverwrite(otherSFM *ShortFileMetadata) bool {
	return sFM.LastChanged.After(otherSFM.LastChanged)
}
