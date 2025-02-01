package repo

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"server/domain"
	"time"
)

type authRepo struct{}

func (a *authRepo) CreateJWT(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	// TODO: move secret key to env
	tokenString, err := token.SignedString([]byte("5asfg67sdftgs57df4g5764sdfg473sd4f62g6sdf3sd2g46sdf352sdf4"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func NewAuthRepo() domain.IAuth {
	return &authRepo{}
}

func (a *authRepo) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func (a *authRepo) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
