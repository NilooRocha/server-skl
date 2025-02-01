package in_memory

import (
	"errors"
	"server/domain"
	"sync"
)

type userRepo struct {
	mu    sync.RWMutex
	users map[string]domain.User
}

func NewUserRepo() *userRepo {
	return &userRepo{
		users: make(map[string]domain.User),
	}
}

func (r *userRepo) Update(user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.users[user.ID]
	if !exists {
		return errors.New("user not found")
	}

	r.users[user.ID] = user

	return nil
}

func (r *userRepo) Create(user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *userRepo) Read(id string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepo) ReadByEmail(email string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return domain.User{}, errors.New("user not found")
}

func (r *userRepo) List() ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	users := make([]domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepo) IsValidUniversityEmail(email string) bool {
	//if !strings.HasSuffix(email, ".edu") {
	//	return false
	//}
	//
	//var rx = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	//if !rx.MatchString(email) {
	//	return false
	//}

	return true
}
