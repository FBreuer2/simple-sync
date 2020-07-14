package db

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"

	"github.com/FBreuer2/simple-sync/lib/sync"
	"golang.org/x/crypto/bcrypt"
)

type MemoryDB struct {
	users                 map[string][]byte
	tokens                map[string][]byte
	shortMetadataStore    map[string]*sync.ShortFileMetadata
	extendedMetadataStore map[string]*sync.ExtendedFileMetadata
	blockStore            map[string][]byte
}

const (
	TOKEN_SIZE = 20
)

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		users:                 make(map[string][]byte),
		tokens:                make(map[string][]byte),
		shortMetadataStore:    make(map[string]*sync.ShortFileMetadata),
		extendedMetadataStore: make(map[string]*sync.ExtendedFileMetadata),
		blockStore:            make(map[string][]byte),
	}
}

func (mDB *MemoryDB) Register(user []byte, password []byte) error {
	if mDB.users[string(user)] != nil {
		return errors.New("User already exists.")
	}

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MaxCost/2)

	if err != nil {
		return err
	}

	mDB.users[string(user)] = hash

	return nil
}

func (mDB *MemoryDB) Login(user []byte, password []byte) error {
	if mDB.users[string(user)] == nil {
		return errors.New("User does not exist.")
	}

	return bcrypt.CompareHashAndPassword(mDB.users[string(user)], password)
}

func (mDB *MemoryDB) Rekey(user []byte, oldPassword []byte, newPassword []byte) error {
	if err := mDB.Login(user, oldPassword); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword(newPassword, bcrypt.MaxCost/2)

	if err != nil {
		return err
	}

	mDB.users[string(user)] = hash

	return nil
}

func (mDB *MemoryDB) GenerateToken(user []byte, password []byte) ([]byte, error) {
	if err := mDB.Login(user, password); err != nil {
		return nil, err
	}

	if token := mDB.tokens[string(user)]; token != nil {
		return token, nil
	}

	newTokenBytes := make([]byte, TOKEN_SIZE)

	readBytes, err := rand.Read(newTokenBytes)

	if err != nil {
		return nil, err
	}

	if readBytes < TOKEN_SIZE {
		return nil, errors.New("Could not read enough random bytes to generate the token.")
	}

	mDB.tokens[string(user)] = newTokenBytes

	return newTokenBytes, nil
}

func (mDB *MemoryDB) ValidateToken(user []byte, token []byte) error {
	if mDB.users[string(user)] == nil {
		return errors.New("User does not exist.")
	}

	if bytes.Equal(token, mDB.tokens[string(user)]) == false {
		return errors.New("Wrong token for this user.")
	}

	return nil
}

func (mDB *MemoryDB) RetrieveShortFileMetadata(user []byte) (*sync.ShortFileMetadata, error) {
	if mDB.users[string(user)] == nil {
		return nil, errors.New("User does not exist.")
	}

	if store := mDB.shortMetadataStore[string(user)]; store == nil {
		return nil, FILE_NOT_AVAILABLE
	} else {
		return store, nil
	}
}

func (mDB *MemoryDB) PutShortFileMetadata(user []byte, metadata *sync.ShortFileMetadata) error {
	mDB.shortMetadataStore[string(user)] = metadata
	return nil
}

func (mDB *MemoryDB) RetrieveExtendedFileMetadata(user []byte) (*sync.ExtendedFileMetadata, error) {
	if mDB.users[string(user)] == nil {
		return nil, errors.New("User does not exist.")
	}

	if store := mDB.extendedMetadataStore[string(user)]; store == nil {
		return nil, FILE_NOT_AVAILABLE
	} else {
		return store, nil
	}
}

func (mDB *MemoryDB) PutExtendedFileMetadata(user []byte, metadata *sync.ExtendedFileMetadata) error {
	mDB.extendedMetadataStore[string(user)] = metadata
	return nil
}

func (mDB *MemoryDB) RetrieveFile(user []byte, file string) (io.Reader, error) {
	eFM, err := mDB.RetrieveExtendedFileMetadata(user)

	if err != nil {
		return nil, err
	}

	blockFile, err := NewBlockFile(eFM, mDB)

	return blockFile, err
}

func (mDB *MemoryDB) HasBlock(hash []byte) bool {
	if block := mDB.blockStore[string(hash)]; block == nil {
		return false
	} else {
		return true
	}
}

func (mDB *MemoryDB) RetrieveBlock(hash []byte) (io.Reader, error) {
	if block := mDB.blockStore[string(hash)]; block == nil {
		return nil, BLOCK_NOT_AVAILABLE
	} else {
		return bytes.NewReader(block), nil
	}
}

func (mDB *MemoryDB) PutBlock(hash []byte, block []byte) error {
	mDB.blockStore[string(hash)] = block
	return nil
}
