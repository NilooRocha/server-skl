package routes

import (
	"log"
	"net/http"
	"server/domain"
	"server/entrypoint/handler"
	"server/entrypoint/middleware"
	usecaseverification "server/usecase/verification"
)

func CreateVerificationRoutes(userRepo domain.IUser, verificationRepo domain.IVerification) {
	log.Println("Initializing verification routes...")

	verificationCodeUseCase := usecaseverification.NewVerifyAccount(userRepo, verificationRepo)
	resendVerification := usecaseverification.NewResendVerification(verificationRepo, userRepo)
	verificationHandler := handler.NewVerificationHandler(verificationCodeUseCase, resendVerification)

	http.Handle("/verify-account", middleware.Cors(http.HandlerFunc(verificationHandler.VerifyAccount)))
	http.Handle("/resend-verification-code", middleware.Cors(http.HandlerFunc(verificationHandler.ResendVerification)))
}
