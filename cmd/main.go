package main

import (
	"github.com/lpernett/godotenv"
	"log"
	"net/http"
	"os"
	"server/infra/repo"
	"server/infra/repo/in_memory"

	"server/infra/routes"
)

func main() {

	log.Println("Initializing server...")

	log.Println("Loading .env file...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	userRepo := in_memory.NewUserRepo()
	idRepo := repo.NewIdRepo()
	verificationRepo := repo.NewVerificationRepo()
	authRepo := repo.NewAuthRepo()

	log.Println("Creating routes")
	routes.CreateRoutes(userRepo, authRepo, idRepo, verificationRepo)

	CreateAdminUserIfNotExists(userRepo, authRepo, idRepo)

	port := os.Getenv("PORT")

	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
