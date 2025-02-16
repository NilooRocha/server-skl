package auth

import (
	"server/domain"
	"server/usecase/_erros"
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
		return LoginOutput{}, errors.ErrInvalidEmailOrPassword
	}

	if !l.auth.VerifyPassword(i.Password, user.Password) {
		return LoginOutput{}, errors.ErrInvalidEmailOrPassword
	}

	accessToken, err := l.auth.CreateAccessToken(user.ID)
	if err != nil {
		return LoginOutput{}, errors.ErrTokenCreationFailed
	}

	refreshToken, err := l.auth.CreateRefreshToken(user.ID)
	if err != nil {
		return LoginOutput{}, errors.ErrTokenCreationFailed
	}

	return LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
