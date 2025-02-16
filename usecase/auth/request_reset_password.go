package auth

import (
	"fmt"
	"server/domain"
	errors "server/usecase/_erros"
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
		return errors.ErrRequestUserNotFound
	}

	resetToken, err := r.authRepo.CreatePasswordResetToken(user.ID)
	if err != nil {
		return errors.ErrTokenCreationFailure
	}

	fmt.Printf("Password reset token for user '%s': %s\n", input.Email, resetToken)
	return nil
}
