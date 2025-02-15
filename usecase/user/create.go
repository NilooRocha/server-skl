package user

import (
	"errors"
	"log"
	"server/domain"
	"time"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrEmailNotValid          = errors.New("email not valid")
	ErrCreatePasswordHash     = errors.New("failed to hash password")
	ErrCreateVerificationCode = errors.New("failed to create verification code")
	ErrCreateId               = errors.New("failed to create id")
	ErrUserCreationFailed     = errors.New("failed to create user")
)

type CreateUserInput struct {
	FullName string
	Email    string
	Password string
}

type CreateUser struct {
	repo         domain.IUser
	auth         domain.IAuth
	id           domain.IId
	verification domain.IVerification
}

func NewCreateUser(
	userRepo domain.IUser,
	authRepo domain.IAuth,
	idRepo domain.IId,
	verificationRepo domain.IVerification,
) *CreateUser {
	return &CreateUser{
		repo:         userRepo,
		auth:         authRepo,
		id:           idRepo,
		verification: verificationRepo,
	}
}

func (uc *CreateUser) Execute(i CreateUserInput) error {
	if userFound, _ := uc.repo.ReadByEmail(i.Email); userFound.ID != "" {
		return ErrEmailAlreadyRegistered
	}

	id, err := uc.id.Create()
	if err != nil {
		return ErrCreateId
	}

	hashedPassword, err := uc.auth.HashPassword(i.Password)
	if err != nil {
		log.Println(err)
		return ErrCreatePasswordHash
	}

	if !uc.repo.IsValidUniversityEmail(i.Email) {
		return ErrEmailNotValid
	}

	user := domain.User{
		ID:         id.Value,
		FullName:   i.FullName,
		Email:      i.Email,
		Password:   hashedPassword,
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Time{},
	}

	if err = uc.repo.Create(user); err != nil {
		log.Println(err)
		return ErrUserCreationFailed
	}

	code, err := uc.verification.GenerateCode()
	if err != nil {
		return ErrCreateVerificationCode
	}

	verification := domain.Verification{
		UserID:     user.ID,
		Code:       code,
		ExpiresAt:  time.Now().Add(5 * time.Minute),
		LastSentAt: time.Now(),
	}

	if err = uc.verification.Store(verification); err != nil {
		return err
	}

	if err = uc.verification.SendVerification(user.ID, code); err != nil {
		return err
	}

	return nil
}
