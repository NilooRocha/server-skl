package auth

import (
	"errors"
	"server/domain"
)

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrTokenCreationFailed    = errors.New("failed to create token")
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
	User  domain.User
}

type Login struct {
	user domain.IUser
	auth domain.IAuth
}

func NewLogin(userRepo domain.IUser, authRepo domain.IAuth) *Login {
	return &Login{
		user: userRepo,
		auth: authRepo,
	}
}

func (l *Login) Execute(i LoginInput) (LoginOutput, error) {
	user, err := l.user.ReadByEmail(i.Email)
	if err != nil {
		return LoginOutput{}, ErrInvalidEmailOrPassword
	}

	if !l.auth.VerifyPassword(i.Password, user.Password) {
		return LoginOutput{}, ErrInvalidEmailOrPassword
	}

	token, err := l.auth.CreateJWT(user.ID)
	if err != nil {
		return LoginOutput{}, ErrTokenCreationFailed
	}

	return LoginOutput{Token: token, User: user}, nil
}
