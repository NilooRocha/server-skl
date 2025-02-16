package verification

import (
	"log"
	"server/domain"
	errors "server/usecase/_erros"
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
		return errors.ErrInvalidCode
	}

	valid, err := v.verificationRepo.Validate(user.ID, i.Code)
	if err != nil || !valid {
		return errors.ErrInvalidCode
	}

	user.IsVerified = valid

	err = v.userRepo.Update(user)

	if err != nil {
		log.Println("Error updating user:", err)
		return errors.ErrUserUpdateFailed
	}

	err = v.verificationRepo.Delete(user.ID)
	if err != nil {
		return err
	}

	return nil
}
