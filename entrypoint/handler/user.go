package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	usecase "server/usecase/user"
	"strings"
)

type UserHandler struct {
	CreateUserUseCase     *usecase.CreateUser
	ReadUserUseCase       *usecase.ReadUser
	ListUserUseCase       *usecase.ListUsers
	FirstTimeSetupUseCase *usecase.FirstTimeSetup
}

func NewUserHandler(
	create *usecase.CreateUser,
	read *usecase.ReadUser,
	list *usecase.ListUsers,
	firstTimeSetup *usecase.FirstTimeSetup,
) *UserHandler {
	return &UserHandler{
		CreateUserUseCase:     create,
		ReadUserUseCase:       read,
		ListUserUseCase:       list,
		FirstTimeSetupUseCase: firstTimeSetup,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	input := usecase.CreateUserInput{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.FullName == "" || input.Email == "" || input.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err := h.CreateUserUseCase.Execute(input)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyRegistered) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		if errors.Is(err, usecase.ErrEmailNotValid) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) FirstTimeSetup(w http.ResponseWriter, r *http.Request) {
	input := usecase.FirstTimeSetupInput{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Location == "" || input.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err := h.FirstTimeSetupUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) ReadUser(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 3 {
		http.Error(w, "ID parameter is missing", http.StatusBadRequest)
		return
	}
	id := urlParts[2]

	output, err := h.ReadUserUseCase.Execute(usecase.ReadUserInput{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(output.User)
	if err != nil {
		return
	}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	output, err := h.ListUserUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(output.Users)
	if err != nil {
		return
	}
}
