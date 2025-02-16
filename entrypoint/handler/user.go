package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"server/usecase/user"
	"strings"
)

type UserHandler struct {
	CreateUserUseCase     *user.CreateUser
	ReadUserUseCase       *user.ReadUser
	ListUserUseCase       *user.ListUsers
	FirstTimeSetupUseCase *user.FirstTimeSetup
	UpdateUserUseCase     *user.UpdateUser
}

func NewUserHandler(
	create *user.CreateUser,
	read *user.ReadUser,
	list *user.ListUsers,
	firstTimeSetup *user.FirstTimeSetup,
	update *user.UpdateUser,
) *UserHandler {
	return &UserHandler{
		CreateUserUseCase:     create,
		ReadUserUseCase:       read,
		ListUserUseCase:       list,
		FirstTimeSetupUseCase: firstTimeSetup,
		UpdateUserUseCase:     update,
	}
}

func getUserIDFromUrl(r *http.Request) (string, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		return "", errors.New("ID parameter is missing")
	}
	return parts[len(parts)-1], nil
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var input user.CreateUserInput
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
		if errors.Is(err, user.ErrEmailAlreadyRegistered) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, user.ErrEmailNotValid) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Error creating user: %v", err)
		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) FirstTimeSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserIDFromUrl(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input user.FirstTimeSetupInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input.ID = userID
	if input.Location == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.FirstTimeSetupUseCase.Execute(input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) ReadUser(w http.ResponseWriter, userID string) {
	output, err := h.ReadUserUseCase.Execute(user.ReadUserInput{ID: userID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(output.User); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	output, err := h.ListUserUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(output.Users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserIDFromUrl(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input user.UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input.ID = userID

	if err := h.UpdateUserUseCase.Execute(input); err != nil {
		if errors.Is(err, user.ErrForbidden) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) UserRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromUrl(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.ReadUser(w, userID)
	case http.MethodPatch:
		h.UpdateUser(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
