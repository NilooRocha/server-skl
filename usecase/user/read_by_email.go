package user

import (
	"server/domain"
)

type ReadByEmailUserInput struct {
	Email string
}

type ReadByEmailUserOutput struct {
	User domain.User
}

type ReadUserByEmail struct {
	repo domain.IUser
}

func NewReadUserByEmail(repo domain.IUser) *ReadUserByEmail {
	return &ReadUserByEmail{repo: repo}
}

func (r *ReadUserByEmail) Execute(input ReadByEmailUserInput) (ReadByEmailUserOutput, error) {
	user, err := r.repo.ReadByEmail(input.Email)
	if err != nil {
		return ReadByEmailUserOutput{}, err
	}

	return ReadByEmailUserOutput{User: user}, nil
}
