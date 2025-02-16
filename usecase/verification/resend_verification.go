package verification

import (
	"log"
	"server/domain"
	"server/usecase/_erros"
	"time"
)

type ResendVerificationInput struct {
	Email string
}

type ResendVerification struct {
	verificationRepo domain.IVerification
	userRepo         domain.IUser
}

func NewResendVerification(
	verificationRepo domain.IVerification,
	userRepo domain.IUser,
) *ResendVerification {
	return &ResendVerification{
		verificationRepo: verificationRepo,
		userRepo:         userRepo,
	}
}
func (rv *ResendVerification) Execute(i ResendVerificationInput) error {
	user, err := rv.userRepo.ReadByEmail(i.Email)
	if err != nil {
		log.Println("Error fetching user:", err)
		return errors.ErrUserNotFound
	}

	verification, err := rv.verificationRepo.Read(user.ID)
	if err != nil {
		log.Println("Error fetching verification data:", err)
		return err
	}

	now := time.Now()

	const resendInterval = 2*time.Minute + 30*time.Second
	if verification.LastSentAt.Add(resendInterval).After(now) {
		return errors.ErrTimeInterval
	}

	var code string
	if verification.ExpiresAt.After(now) {
		code = verification.Code
	} else {
		code, err = rv.verificationRepo.GenerateCode()
		if err != nil {
			log.Println("Error generating new verification code:", err)
			return errors.ErrGenerateCode
		}

		verification.Code = code
		verification.ExpiresAt = now.Add(5 * time.Minute)
	}

	verification.LastSentAt = now
	if err := rv.verificationRepo.Store(verification); err != nil {
		log.Println("Error storing verification data:", err)
		return err
	}

	return rv.verificationRepo.SendVerification(user.ID, code)
}
