package db

import (
	"bytes"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type MemoryDB struct {
	users 	map[[]byte][]byte
	tokens 	map[[]byte][]byte
}

const (
	TOKEN_SIZE	20
)

func NewMemoryDB() (*MemoryDB) {
	return &MemoryDB{}
}


func (mDB *MemoryDB) Init() (error) {
	mDB.users = make(map[[]byte][]byte)

}


func (mDB *MemoryDB) Register(user []byte, password []byte) (error) {
	if users[user] != nil {
		return errors.New("User already exists.")
	}

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MaxCost/2)

	if (err != nil) {
		return err
	}

	users[user] = hash

	return nil
}


func (mDB *MemoryDB) Login(user []byte, password []byte) (error) {
	if users[user] == nil {
		return errors.New("User does not exist.")
	}

    return bcrypt.CompareHashAndPassword(users[user], password)
}


func (mDB *MemoryDB) Rekey(user []byte, oldPassword []byte, newPassword []byte) (error) {
	if err := mDB.Login(user, oldPassword); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword(newPassword, bcrypt.MaxCost/2)

	if (err != nil) {
		return err
	}

	users[user] = hash

	return nil
}


func (mDB *MemoryDB) GetToken(user []byte, password []byte) ([]byte, error) {
	if (err := mDB.Login(user, password); err != nil) {
		return nil, err
	}

	if token := mDB.tokens[user]; token != nil {
		return token, nil
	}

	newTokenBytes := make([]byte, TOKEN_SIZE)

	readBytes, err := rand.Read(newTokenBytes) (n int, err error)

	if (err != nil) {
		return nil, err
	}

	if (readBytes < TOKEN_SIZE) {
		return nil, erros.New("Could not read enough random bytes to generate the token.")
	}

	mDB.tokens[user] = newTokenBytes

	return newTokenBytes, nil
}


func (mDB *MemoryDB) ValidateToken(user []byte, token []byte) (error) {
	if users[user] == nil {
		return errors.New("User does not exist.")
	}

	if bytes.Equal(token, users[user]) == false {
		return errors.New("Wrong token for this user.")
	}

	return nil
}