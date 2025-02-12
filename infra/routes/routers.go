package routes

import (
	"log"
	"net/http"
	"server/entrypoint/handler"
	"server/entrypoint/middleware"
	"server/infra/repo"
	"server/infra/repo/in_memory"
	usecaseauth "server/usecase/auth"
	usecaseuser "server/usecase/user"
	usecaseverification "server/usecase/verification"
)

func CreateRoutes() {
	log.Println("Initializing user routes...")

	userRepo := in_memory.NewUserRepo()
	idRepo := repo.NewIdRepo()
	verificationRepo := repo.NewVerificationRepo()
	authRepo := repo.NewAuthRepo()

	//---------------VERIFICATION----------------------
	verificationCodeUseCase := usecaseverification.NewVerifyAccount(userRepo, verificationRepo)
	resendVerification := usecaseverification.NewResendVerification(verificationRepo, userRepo)
	verificationHandler := handler.NewVerificationHandler(verificationCodeUseCase, resendVerification)

	http.Handle("/verify-account", middleware.Cors(http.HandlerFunc(verificationHandler.VerifyAccount)))
	http.Handle("/resend-verification-code", middleware.Cors(http.HandlerFunc(verificationHandler.ResendVerification)))

	//---------------USER----------------------

	createUserUseCase := usecaseuser.NewCreateUser(userRepo, authRepo, idRepo, verificationRepo)
	readUserUseCase := usecaseuser.NewReadUser(userRepo)
	listUsersUseCase := usecaseuser.NewListUsers(userRepo)
	firstTimeSetupUseCase := usecaseuser.NewFirstTimeSetup(userRepo)
	updateLocationUseCase := usecaseuser.NewUpdateLocation(userRepo)

	userHandler := handler.NewUserHandler(createUserUseCase, readUserUseCase, listUsersUseCase, firstTimeSetupUseCase, updateLocationUseCase)

	http.Handle("/users", middleware.Cors(middleware.RequireAuth(http.HandlerFunc(userHandler.ListUsers), userRepo, authRepo)))
	http.Handle("/user", middleware.Cors(http.HandlerFunc(userHandler.CreateUser)))
	http.Handle("/user/", middleware.Cors(middleware.RequireAuth(http.HandlerFunc(userHandler.ReadUser), userRepo, authRepo)))
	http.Handle("/user/first-time-setup/", middleware.Cors(middleware.RequireAuth(http.HandlerFunc(userHandler.FirstTimeSetup), userRepo, authRepo)))
	http.Handle("/user/update-location/", middleware.Cors(middleware.RequireAuth(http.HandlerFunc(userHandler.UpdateLocation), userRepo, authRepo)))

	//---------------AUTH----------------------

	loginUseCase := usecaseauth.NewLogin(userRepo, authRepo)
	resetPasswordUseCase := usecaseauth.NewResetPassword(userRepo, authRepo)
	requestResetPasswordUseCase := usecaseauth.NewRequestResetPassword(userRepo, authRepo)
	authHandler := handler.NewAuthHandler(loginUseCase, resetPasswordUseCase, requestResetPasswordUseCase)

	http.Handle("/login", middleware.Cors(http.HandlerFunc(authHandler.Login)))
	http.Handle("/logout", middleware.Cors(middleware.RequireAuth(http.HandlerFunc(authHandler.Logout), userRepo, authRepo)))
	http.Handle("/reset-password", middleware.Cors(http.HandlerFunc(authHandler.ResetPassword)))
	http.Handle("/request-reset-password", middleware.Cors(http.HandlerFunc(authHandler.RequestResetPassword)))
}
