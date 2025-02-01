package domain

import "time"

type User struct {
	ID             string    `json:"id"`
	FullName       string    `json:"fullName"`
	Email          string    `json:"email"`
	Password       string    `json:"-"`
	UniversityName string    `json:"universityName"`
	Location       string    `json:"location"`
	IsVerified     bool      `json:"isVerified"`
	Balance        float64   ` json:"balance"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type IUser interface {
	Create(user User) error
	Read(id string) (User, error)
	Update(user User) error
	ReadByEmail(id string) (User, error)
	List() ([]User, error)
	IsValidUniversityEmail(string) bool
}
