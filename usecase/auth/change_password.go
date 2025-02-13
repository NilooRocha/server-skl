package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"server/domain"
	"time"
)

var (
	ErrPasswordUpdateFailed = errors.New("failed to update password")
	ErrPasswordIncorrect    = errors.New("incorrect old password")
)

type ChangePasswordInput struct {
	ID              string
	CurrentPassword string
	NewPassword     string
}

type ChangePassword struct {
	repo domain.IUser
}

func NewChangePassword(userRepo domain.IUser) *ChangePassword {
	return &ChangePassword{
		repo: userRepo,
	}
}

func (u *ChangePassword) Execute(input ChangePasswordInput) error {
	user, err := u.repo.Read(input.ID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		return ErrPasswordIncorrect
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return ErrPasswordUpdateFailed
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := u.repo.Update(user); err != nil {
		return ErrPasswordUpdateFailed
	}

	return nil
}
