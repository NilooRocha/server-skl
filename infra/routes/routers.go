package routes

import (
	"log"
	"server/domain"
)

func CreateRoutes(userRepo domain.IUser, authRepo domain.IAuth, idRepo domain.IId, verificationRepo domain.IVerification) {
	log.Println("Initializing all routes...")

	CreateVerificationRoutes(userRepo, verificationRepo)
	CreateUserRoutes(userRepo, authRepo, idRepo, verificationRepo)
	CreateAuthRoutes(userRepo, authRepo)

}
