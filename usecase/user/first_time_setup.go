package user

import (
	"errors"
	"log"
	"server/domain"
	"time"
)

var (
	ErrUserUpdateFailed = errors.New("failed to update user")
)

type FirstTimeSetupInput struct {
	Email    string
	Location string
}

type FirstTimeSetup struct {
	repo domain.IUser
}

func NewFirstTimeSetup(userRepo domain.IUser) *FirstTimeSetup {
	return &FirstTimeSetup{
		repo: userRepo,
	}
}

func (fts *FirstTimeSetup) Execute(i FirstTimeSetupInput) error {
	user, err := fts.repo.ReadByEmail(i.Email)
	if err != nil {
		return err
	}

	user.Location = i.Location
	user.UpdatedAt = time.Now()

	err = fts.repo.Update(user)

	if err != nil {
		log.Println("Error updating user:", err)
		return ErrUserUpdateFailed
	}

	return nil
}
