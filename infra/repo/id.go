// in_memory/id_repo.go

package repo

import (
	"math/rand"
	"server/domain"
	"time"
)

type idRepo struct{}

func NewIdRepo() domain.IId {
	return &idRepo{}
}

func (r *idRepo) Create() (domain.Id, error) {
	id := domain.Id{Value: GenerateShortID()}
	return id, nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const length = 8

func GenerateShortID() string {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()),
	)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]

	}
	return string(b)
}
