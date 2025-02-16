package user

import (
	"log"
	"server/domain"
	errors "server/usecase/_erros"
	"time"
)

type FirstTimeSetupInput struct {
	ID       string
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
	user, err := fts.repo.Read(i.ID)
	if err != nil {
		return err
	}

	user.Location = i.Location
	user.UpdatedAt = time.Now()

	err = fts.repo.Update(user)

	if err != nil {
		log.Println("Error updating user:", err)
		return errors.ErrUserUpdateFailed
	}

	return nil
}
