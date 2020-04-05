package db

type AuthenticatorDatabase interface {
	Register(user string, password string) (error)
	Login(user string, password string) (error)

	GetToken(user string, password string) (string, error)
	ValidateToken(user string, token string) (error)
}