package user

import (
	"server/domain"
)

type ReadUserInput struct {
	ID string
}

type ReadUserOutput struct {
	User domain.User
}

type ReadUser struct {
	repo domain.IUser
}

func NewReadUser(repo domain.IUser) *ReadUser {
	return &ReadUser{repo: repo}
}

func (ru *ReadUser) Execute(input ReadUserInput) (ReadUserOutput, error) {
	user, err := ru.repo.Read(input.ID)
	if err != nil {
		return ReadUserOutput{}, err
	}

	return ReadUserOutput{User: user}, nil
}
