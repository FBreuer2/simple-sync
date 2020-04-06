package db

type AuthenticatorDatabase interface {
	Init() (error)
	Register(user []byte, password []byte) (error)
	Login(user []byte, password []byte) (error)
	Rekey(user []byte, oldPassword []byte, newPassword []byte) (error)

	GenerateToken(user []byte, password []byte) ([]byte, error)
	ValidateToken(user []byte, token []byte) (error)
}