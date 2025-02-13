package handler

import (
	"encoding/json"
	"net/http"
	usecase "server/usecase/auth"
	"time"
)

type AuthHandler struct {
	LoginUseCase                *usecase.Login
	ResetPasswordUseCase        *usecase.ResetPassword
	RequestResetPasswordUseCase *usecase.RequestResetPassword
	ChangePasswordUseCase       *usecase.ChangePassword
}

func NewAuthHandler(login *usecase.Login, resetPassword *usecase.ResetPassword, requestResetPassword *usecase.RequestResetPassword, changePassword *usecase.ChangePassword) *AuthHandler {
	return &AuthHandler{
		LoginUseCase:                login,
		ResetPasswordUseCase:        resetPassword,
		RequestResetPasswordUseCase: requestResetPassword,
		ChangePasswordUseCase:       changePassword,
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

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input usecase.ResetPasswordInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if input.ResetToken == "" || input.NewPassword == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.ResetPasswordUseCase.Execute(input); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) RequestResetPassword(w http.ResponseWriter, r *http.Request) {

	var input usecase.RequestResetPasswordInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if input.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if err := h.RequestResetPasswordUseCase.Execute(input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset link sent"))
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPatch {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserIDFromUrl(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var input usecase.ChangePasswordInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if input.CurrentPassword == "" || input.NewPassword == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	input.ID = userID

	if err := h.ChangePasswordUseCase.Execute(input); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
