package user

import (
	"errors"
	"server/domain"
	"server/permissions"
	"time"
)

var (
	ErrUpdateFailed = errors.New("failed to update user")
	ErrForbidden    = errors.New("permission denied")
)

type UpdateInput struct {
	ID         string `json:"id"`
	FullName   string `json:"fullName"`
	Email      string `json:"email"`
	Location   string `json:"location"`
	IsVerified *bool  `json:"isVerified"`
}

type UpdateUser struct {
	repo domain.IUser
}

func NewUpdateUser(userRepo domain.IUser) *UpdateUser {
	return &UpdateUser{
		repo: userRepo,
	}
}

func (u *UpdateUser) Execute(i UpdateInput) error {
	user, err := u.repo.Read(i.ID)
	if err != nil {
		return err
	}

	if !permissions.Can(user.Role, "update", "user") {
		return ErrForbidden
	}

	if user.Role.Name == domain.UserRole && (i.Email != "" || i.IsVerified != nil) {
		return ErrForbidden
	}

	if i.Location != "" {
		user.Location = i.Location
	}

	if i.FullName != "" {
		user.FullName = i.FullName
	}

	if i.Email != "" {
		user.Email = i.Email
	}

	if i.IsVerified != nil {
		user.IsVerified = *i.IsVerified
	}

	user.UpdatedAt = time.Now()

	if err := u.repo.Update(user); err != nil {
		return ErrUpdateFailed
	}

	return nil
}
