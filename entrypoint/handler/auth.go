package handler

import (
	"encoding/json"
	"net/http"
	usecase "server/usecase/auth"
	"time"
)

type AuthHandler struct {
	LoginUseCase *usecase.Login
}

func NewAuthHandler(login *usecase.Login) *AuthHandler {
	return &AuthHandler{
		LoginUseCase: login,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	input := usecase.LoginInput{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	output, err := h.LoginUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    output.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: set secure to true
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(15 * time.Minute),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "RefreshToken",
		Value:    output.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: set secure to true
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output.User)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: Defina para true em produção
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "RefreshToken",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: Defina para true em produção
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
	})

	w.Header().Set("Cache-Control", "no-store")

	w.WriteHeader(http.StatusOK)
}
