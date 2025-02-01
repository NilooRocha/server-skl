package handler

import (
	"encoding/json"
	"net/http"
	usecase "server/usecase/verification"
)

type VerificationHandler struct {
	VerifyAccountUseCase      *usecase.VerifyAccount
	ResendVerificationUseCase *usecase.ResendVerification
}

func NewVerificationHandler(
	verifyAccount *usecase.VerifyAccount,
	resendVerification *usecase.ResendVerification,
) *VerificationHandler {
	return &VerificationHandler{
		VerifyAccountUseCase:      verifyAccount,
		ResendVerificationUseCase: resendVerification,
	}
}

func (v *VerificationHandler) VerifyAccount(w http.ResponseWriter, r *http.Request) {
	input := usecase.VerifyAccountInput{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Code == "" {
		http.Error(w, "Missing email or verification code", http.StatusBadRequest)
		return
	}

	if err := v.VerifyAccountUseCase.Execute(input); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (v *VerificationHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	input := usecase.ResendVerificationInput{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Email == "" {
		http.Error(w, "Missing email", http.StatusBadRequest)
		return
	}

	if err := v.ResendVerificationUseCase.Execute(input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
