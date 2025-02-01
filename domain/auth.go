package domain

type IAuth interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) bool
	CreateAccessToken(userID string) (string, error)
	CreateRefreshToken(userID string) (string, error)
	ValidateRefreshToken(tokenString string) (string, error)
}
