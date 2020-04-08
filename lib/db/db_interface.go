package db

import "github.com/FBreuer2/simple-sync/lib/sync"

type AuthenticatorDatabase interface {
	Register(user []byte, password []byte) error
	Login(user []byte, password []byte) error
	Rekey(user []byte, oldPassword []byte, newPassword []byte) error

	GenerateToken(user []byte, password []byte) ([]byte, error)
	ValidateToken(user []byte, token []byte) error
}

type FileMetadataDatabase interface {
	RetrieveShortFileMetadata(user []byte) (*sync.ShortFileMetadata, error)
	PutShortFileMetadata(user []byte, metadata *sync.ShortFileMetadata) error

	RetrieveExtendedFileMetadata(user []byte) (*sync.ExtendedFileMetadata, error)
	PutExtendedFileMetadata(user []byte, metadata *sync.ExtendedFileMetadata) error
}

type BlockDatabase interface {
	RetrieveBlock(hash []byte) ([]byte, error)
	PutBlock(hash []byte, block []byte) error
}

type FullDatabase interface {
	AuthenticatorDatabase
	FileMetadataDatabase
	BlockDatabase
}
