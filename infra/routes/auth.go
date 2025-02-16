package routes

import (
	"log"
	"net/http"
	"server/domain"
	"server/entrypoint/handler"
	"server/entrypoint/middleware"
	usecaseauth "server/usecase/auth"
)

func CreateAuthRoutes(userRepo domain.IUser, authRepo domain.IAuth) {
	log.Println("Initializing auth routes...")

	loginUseCase := usecaseauth.NewLogin(userRepo, authRepo)
	resetPasswordUseCase := usecaseauth.NewResetPassword(userRepo, authRepo)
	requestResetPasswordUseCase := usecaseauth.NewRequestResetPassword(userRepo, authRepo)
	changePasswordUseCase := usecaseauth.NewChangePassword(userRepo)

	authHandler := handler.NewAuthHandler(loginUseCase, resetPasswordUseCase, requestResetPasswordUseCase, changePasswordUseCase)

	http.Handle("/login", middleware.Cors(http.HandlerFunc(authHandler.Login)))
	http.Handle("/logout", middleware.Cors(middleware.RequireAuth(http.HandlerFunc(authHandler.Logout), userRepo, authRepo)))
	http.Handle("/change-password/", middleware.Cors(middleware.RequireAuth(http.HandlerFunc(authHandler.ChangePassword), userRepo, authRepo)))
	http.Handle("/reset-password", middleware.Cors(http.HandlerFunc(authHandler.ResetPassword)))
	http.Handle("/request-reset-password", middleware.Cors(http.HandlerFunc(authHandler.RequestResetPassword)))
}
