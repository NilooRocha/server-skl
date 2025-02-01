package repo

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"server/domain"
	"time"
)

type authRepo struct{}

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

func (a *authRepo) CreateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Minute * 15).Unix(), // 15 minutos
	})

	// TODO: move secret key to env and create a new one
	tokenString, err := token.SignedString([]byte("5asfg67sdftgs57df4g5764sdfg473sd4f62g6sdf3sd2g46sdf352sdf4"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *authRepo) CreateRefreshToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 dias
	})

	// TODO: move secret key to env
	tokenString, err := token.SignedString([]byte("5asfg67sdftgs57df4g5764sdfg473sd4f62g6sdf3sd2g46sdf352sdf4"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *authRepo) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("5asfg67sdftgs57df4g5764sdfg473sd4f62g6sdf3sd2g46sdf352sdf4"), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	}

	return "", fmt.Errorf("invalid token")
}
