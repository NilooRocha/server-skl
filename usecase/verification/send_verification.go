package verification

import (
	"server/domain"
)

type SendVerificationInput struct {
	UserID string
	Code   string
}

type SendVerification struct {
	VerificationData domain.IVerification
}

func NewSendVerification(verificationDataRepo domain.IVerification) *SendVerification {
	return &SendVerification{
		VerificationData: verificationDataRepo,
	}
}

func (gv *SendVerification) Execute(i SendVerificationInput) error {
	err := gv.VerificationData.SendVerification(i.UserID, i.Code)
	if err != nil {
		return err
	}

	return nil
}
