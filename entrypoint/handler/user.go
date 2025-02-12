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
	UpdateLocationUseCase *usecase.UpdateLocation
}

func NewUserHandler(
	create *usecase.CreateUser,
	read *usecase.ReadUser,
	list *usecase.ListUsers,
	firstTimeSetup *usecase.FirstTimeSetup,
	updateLocation *usecase.UpdateLocation,
) *UserHandler {
	return &UserHandler{
		CreateUserUseCase:     create,
		ReadUserUseCase:       read,
		ListUserUseCase:       list,
		FirstTimeSetupUseCase: firstTimeSetup,
		UpdateLocationUseCase: updateLocation,
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

	path := r.URL.Path
	parts := strings.Split(path, "/")
	userId := parts[len(parts)-1]

	input := usecase.FirstTimeSetupInput{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input.ID = userId

	if input.Location == "" {
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
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id := parts[len(parts)-1]
	if len(parts) < 3 {
		http.Error(w, "ID parameter is missing", http.StatusBadRequest)
		return
	}

	output, err := h.ReadUserUseCase.Execute(usecase.ReadUserInput{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(output.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func (h *UserHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	parts := strings.Split(path, "/")
	userId := parts[len(parts)-1]

	input := usecase.UpdateLocationInput{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input.ID = userId

	if input.Location == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err := h.UpdateLocationUseCase.Execute(input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
