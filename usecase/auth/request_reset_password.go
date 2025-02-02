package auth

import (
	"errors"
	"fmt"
	"server/domain"
)

var (
	ErrRequestUserNotFound  = errors.New("user not found")
	ErrTokenCreationFailure = errors.New("failed to create reset token")
)

type RequestResetPassword struct {
	userRepo domain.IUser
	authRepo domain.IAuth
}

type RequestResetPasswordInput struct {
	Email string
}

func NewRequestResetPassword(userRepo domain.IUser, authRepo domain.IAuth) *RequestResetPassword {
	return &RequestResetPassword{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (r *RequestResetPassword) Execute(input RequestResetPasswordInput) error {
	user, err := r.userRepo.ReadByEmail(input.Email)
	if err != nil {
		return ErrRequestUserNotFound
	}

	resetToken, err := r.authRepo.CreatePasswordResetToken(user.ID)
	if err != nil {
		return ErrTokenCreationFailure
	}

	fmt.Printf("Password reset token for user '%s': %s\n", input.Email, resetToken)
	return nil
}
