package verification

import (
	"errors"
	"log"
	"server/domain"
	"time"
)

var (
	ErrErrorWhenStore = errors.New("failed to store verification data")
)

type GenerateInput struct {
	UserID string
}

type GenerateVerification struct {
	repo domain.IVerification
}

func NewGenerateVerification(verificationRepo domain.IVerification) *GenerateVerification {
	return &GenerateVerification{
		repo: verificationRepo,
	}
}

func (gv *GenerateVerification) Execute(i GenerateInput) (string, error) {
	code, err := gv.repo.GenerateCode()
	if err != nil {
		return "", err
	}

	ttl := 5 * time.Minute

	VerificationData := domain.Verification{
		UserID:     i.UserID,
		Code:       code,
		ExpiresAt:  time.Now().Add(ttl),
		LastSentAt: time.Now(),
	}

	err = gv.repo.Store(VerificationData)
	if err != nil {
		log.Println("Error storing user verification code:", err)
		return "", ErrErrorWhenStore
	}

	return code, nil
}
