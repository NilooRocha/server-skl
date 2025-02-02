package auth

import (
	"errors"
	"fmt"
	"server/domain"
)

var (
	ErrInvalidResetToken = errors.New("invalid or expired reset token")
	ErrUserNotFound      = errors.New("user not found")
)

type ResetPasswordInput struct {
	ResetToken  string
	NewPassword string
}

type ResetPassword struct {
	userRepo domain.IUser
	authRepo domain.IAuth
}

func NewResetPassword(userRepo domain.IUser, authRepo domain.IAuth) *ResetPassword {
	return &ResetPassword{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (r *ResetPassword) Execute(input ResetPasswordInput) error {
	userID, err := r.authRepo.ValidatePasswordResetToken(input.ResetToken)
	if err != nil {
		return ErrInvalidResetToken
	}

	user, err := r.userRepo.Read(userID)
	if err != nil {
		return ErrUserNotFound
	}

	hashedPassword, err := r.authRepo.HashPassword(input.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	return r.userRepo.Update(user)
}
