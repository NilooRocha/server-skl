package main

import (
	"log"
	"net/http"

	"server/infra/routes"
)

func main() {

	log.Println("Initializing server...")

	log.Println("Creating routes")
	routes.CreateRoutes()

	port := "8080"

	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
