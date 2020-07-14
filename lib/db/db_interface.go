package db

import (
	"errors"
	"io"

	"github.com/FBreuer2/simple-sync/lib/sync"
)

type AuthenticatorDatabase interface {
	Register(user []byte, password []byte) error
	Login(user []byte, password []byte) error
	Rekey(user []byte, oldPassword []byte, newPassword []byte) error

	GenerateToken(user []byte, password []byte) ([]byte, error)
	ValidateToken(user []byte, token []byte) error
}

type FileDatabase interface {
	RetrieveShortFileMetadata(user []byte) (*sync.ShortFileMetadata, error)
	PutShortFileMetadata(user []byte, metadata *sync.ShortFileMetadata) error

	RetrieveExtendedFileMetadata(user []byte) (*sync.ExtendedFileMetadata, error)
	PutExtendedFileMetadata(user []byte, metadata *sync.ExtendedFileMetadata) error

	RetrieveFile(user []byte, file string) (io.Reader, error)
}

type BlockDatabase interface {
	HasBlock(hash []byte) bool
	RetrieveBlock(hash []byte) (io.Reader, error)
	PutBlock(hash []byte, block []byte) error
}

type FullDatabase interface {
	AuthenticatorDatabase
	FileDatabase
	BlockDatabase
}

var FILE_NOT_AVAILABLE = errors.New("File is not available.")
var BLOCK_NOT_AVAILABLE = errors.New("Block is not available.")

func NewBlockFile(eFM *sync.ExtendedFileMetadata, blockStorage BlockDatabase) (io.Reader, error) {
	readers := make([]io.Reader, len(eFM.StrongBlockHashes))
	for index, strongHash := range eFM.StrongBlockHashes {
		reader, err := blockStorage.RetrieveBlock(strongHash)

		if err != nil {
			return nil, err
		}

		readers[index] = reader
	}

	return io.MultiReader(readers...), nil
}
