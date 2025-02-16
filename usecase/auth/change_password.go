package auth

import (
	"golang.org/x/crypto/bcrypt"
	"server/domain"
	"server/permissions"
	errors "server/usecase/_erros"
	"time"
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

	if !permissions.Can(user.Role, "update", "user") {
		return errors.ErrForbidden
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		return errors.ErrPasswordIncorrect
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrPasswordUpdateFailed
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := u.repo.Update(user); err != nil {
		return errors.ErrPasswordUpdateFailed
	}

	return nil
}
