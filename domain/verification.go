package domain

import "time"

type Verification struct {
	UserID     string
	Code       string
	ExpiresAt  time.Time
	LastSentAt time.Time
}

type IVerification interface {
	GenerateCode() (string, error)
	Read(userID string) (Verification, error)
	Store(verification Verification) error
	Validate(userID, code string) (bool, error)
	Delete(userID string) error
	SendVerification(userID string, code string) error
}
