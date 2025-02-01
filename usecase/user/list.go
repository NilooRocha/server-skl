package user

import (
	"server/domain"
)

type ListUsersOutput struct {
	Users []domain.User
}

type ListUsers struct {
	repo domain.IUser
}

func NewListUsers(repo domain.IUser) *ListUsers {
	return &ListUsers{repo: repo}
}

func (lu *ListUsers) Execute() (ListUsersOutput, error) {
	users, err := lu.repo.List()
	if err != nil {
		return ListUsersOutput{}, err
	}

	return ListUsersOutput{Users: users}, nil
}
