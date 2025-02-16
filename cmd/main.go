package main

import (
	"log"
	"net/http"
	"server/infra/repo"
	"server/infra/repo/in_memory"

	"server/infra/routes"
)

func main() {

	log.Println("Initializing server...")

	userRepo := in_memory.NewUserRepo()
	idRepo := repo.NewIdRepo()
	verificationRepo := repo.NewVerificationRepo()
	authRepo := repo.NewAuthRepo()

	log.Println("Creating routes")
	routes.CreateRoutes(userRepo, authRepo, idRepo, verificationRepo)

	port := "8080"

	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
