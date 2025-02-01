package domain

type IAuth interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) bool
	CreateJWT(email string) (string, error)
}
