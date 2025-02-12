package user

import (
	"errors"
	"log"
	"server/domain"
	"time"
)

var (
	ErrUpdateFailed = errors.New("failed to update user")
)

type UpdateLocationInput struct {
	ID       string
	Location string
}

type UpdateLocation struct {
	repo domain.IUser
}

func NewUpdateLocation(userRepo domain.IUser) *UpdateLocation {
	return &UpdateLocation{
		repo: userRepo,
	}
}

func (fts *UpdateLocation) Execute(i UpdateLocationInput) error {
	user, err := fts.repo.Read(i.ID)
	if err != nil {
		return err
	}

	user.Location = i.Location
	user.UpdatedAt = time.Now()

	err = fts.repo.Update(user)

	if err != nil {
		log.Println("Error updating user:", err)
		return ErrUpdateFailed
	}

	return nil
}
