package repo

import (
	"errors"
	"math/rand"
	"server/domain"
	"sync"
	"time"
)

const (
	sequence   = "1234567890"
	lengthCode = 4
)

type verificationRepo struct {
	mu            sync.RWMutex
	verifications map[string]domain.Verification
	expiration    time.Duration
}

func NewVerificationRepo() *verificationRepo {
	return &verificationRepo{
		verifications: make(map[string]domain.Verification),
		expiration:    5 * time.Minute,
	}
}

func (r *verificationRepo) Store(verificationData domain.Verification) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.verifications[verificationData.UserID] = verificationData

	return nil
}

func (r *verificationRepo) GenerateCode() (string, error) {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()),
	)

	b := make([]byte, lengthCode)
	for i := range b {
		b[i] = sequence[seededRand.Intn(len(sequence))]
	}

	return string(b), nil
}

func (r *verificationRepo) Validate(userID, code string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	verificationData, found := r.verifications[userID]
	if !found {
		return false, errors.New("verification code not found")
	}

	if verificationData.Code != code {
		return false, errors.New("invalid verification code")
	}

	if verificationData.ExpiresAt.Before(time.Now()) {
		return false, errors.New("verification code expired")
	}

	return true, nil
}

func (r *verificationRepo) Delete(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.verifications[userID]; !found {
		return errors.New("verification code not found")
	}

	delete(r.verifications, userID)
	return nil
}

func (r *verificationRepo) SendVerification(userID, code string) error {
	println("Sending verification code:", code, "to user:", userID)
	return nil
}

func (r *verificationRepo) Read(userID string) (domain.Verification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	verification, exists := r.verifications[userID]
	if !exists {
		return domain.Verification{}, errors.New("verification data not found")
	}

	return verification, nil
}
