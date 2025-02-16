package in_memory

import (
	"errors"
	"log"
	"server/domain"
	"sync"
	"time"
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

func (r *userRepo) CreateAdminIfNotExists(auth domain.IAuth, id domain.IId) error {
	// Bloqueio de escrita para garantir atomicidade
	r.mu.Lock()
	defer r.mu.Unlock()

	adminEmail := "admin@domain.com"

	// Verifique se o admin já existe
	userFound, err := r.ReadByEmail(adminEmail)
	if err == nil && userFound.ID != "" {
		// Se já existe, não faz nada
		log.Println("Admin user already exists.")
		return nil
	}

	// Caso não exista, crie o admin
	adminID, err := id.Create()
	if err != nil {
		log.Println("Failed to generate admin ID:", err)
		return err
	}

	hashedPassword, err := auth.HashPassword("admin123")
	if err != nil {
		log.Println("Failed to hash password:", err)
		return err
	}

	// Criar o usuário admin
	adminUser := domain.User{
		ID:         adminID.Value,
		FullName:   "Administrator",
		Email:      adminEmail,
		Location:   "Heaven",
		Password:   hashedPassword,
		Role:       domain.Roles[domain.AdminRole],
		IsVerified: true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Time{},
	}

	// Salvar o novo admin
	if err := r.Create(adminUser); err != nil {
		log.Println("Failed to create admin user:", err)
		return err
	}

	log.Println("Admin user created successfully.")
	return nil
}
