package routes

import (
	"log"
	"net/http"
	"server/domain"
	"server/entrypoint/handler"
	"server/entrypoint/middleware"
	usecaseuser "server/usecase/user"
)

func CreateUserRoutes(userRepo domain.IUser, authRepo domain.IAuth, idRepo domain.IId, verificationRepo domain.IVerification) {
	log.Println("Initializing user routes...")

	createUserUseCase := usecaseuser.NewCreateUser(userRepo, authRepo, idRepo, verificationRepo)
	readUserUseCase := usecaseuser.NewReadUser(userRepo)
	listUsersUseCase := usecaseuser.NewListUsers(userRepo)
	firstTimeSetupUseCase := usecaseuser.NewFirstTimeSetup(userRepo)
	updateLocationUseCase := usecaseuser.NewUpdateUser(userRepo)

	userHandler := handler.NewUserHandler(createUserUseCase, readUserUseCase, listUsersUseCase, firstTimeSetupUseCase, updateLocationUseCase)

	http.Handle("/users", middleware.Cors(middleware.RequireAuth(http.HandlerFunc(userHandler.ListUsers), userRepo, authRepo)))
	http.Handle("/user", middleware.Cors(http.HandlerFunc(userHandler.CreateUser)))
	http.Handle("/user/", middleware.Cors(middleware.RequireAuth(http.HandlerFunc(userHandler.UserRequestHandler), userRepo, authRepo)))
	http.Handle("/user/first-time-setup/", middleware.Cors(middleware.RequireAuth(http.HandlerFunc(userHandler.FirstTimeSetup), userRepo, authRepo)))
}
