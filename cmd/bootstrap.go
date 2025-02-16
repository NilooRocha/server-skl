package main

import (
	"log"
	"server/domain"
)

func CreateAdminUserIfNotExists(repo domain.IUser, auth domain.IAuth, id domain.IId) {
	err := repo.CreateAdminIfNotExists(auth, id)
	if err != nil {
		log.Fatalf("Error creating admin user: %v", err)
	}
}
