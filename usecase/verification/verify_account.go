package verification

import (
	"errors"
	"log"
	"server/domain"
)

var (
	ErrInvalidCode      = errors.New("invalid email or verification code")
	ErrUserUpdateFailed = errors.New("failed to update user")
)

type VerifyAccountInput struct {
	Code  string
	Email string
}

type VerifyAccount struct {
	userRepo         domain.IUser
	verificationRepo domain.IVerification
}

func NewVerifyAccount(userRepo domain.IUser, verificationRepo domain.IVerification) *VerifyAccount {
	return &VerifyAccount{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
	}
}

func (v *VerifyAccount) Execute(i VerifyAccountInput) error {
	user, err := v.userRepo.ReadByEmail(i.Email)
	if err != nil {
		return ErrInvalidCode
	}

	valid, err := v.verificationRepo.Validate(user.ID, i.Code)
	if err != nil || !valid {
		return ErrInvalidCode
	}

	user.IsVerified = valid

	err = v.userRepo.Update(user)

	if err != nil {
		log.Println("Error updating user:", err)
		return ErrUserUpdateFailed
	}

	err = v.verificationRepo.Delete(user.ID)
	if err != nil {
		return err
	}

	return nil
}
