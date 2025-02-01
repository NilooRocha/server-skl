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
	AccessToken  string
	RefreshToken string
	User         domain.User
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

	accessToken, err := l.auth.CreateAccessToken(user.ID)
	if err != nil {
		return LoginOutput{}, ErrTokenCreationFailed
	}

	refreshToken, err := l.auth.CreateRefreshToken(user.ID)
	if err != nil {
		return LoginOutput{}, ErrTokenCreationFailed
	}

	return LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
