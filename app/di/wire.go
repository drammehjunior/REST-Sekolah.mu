package di

import (
	"exampleclean.com/refactor/app/api"
	"exampleclean.com/refactor/app/api/handler"
	"exampleclean.com/refactor/app/config"
	"exampleclean.com/refactor/app/db"
	"exampleclean.com/refactor/app/repository"
	"exampleclean.com/refactor/app/usecase"
)

func InitializeAPI(cfg config.Config) (*api.ServerHTTP, error) {
	gormDB, err := db.ConnectDatabaseGorm(cfg)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepositoryGorm(gormDB)
	userUseCase := usecase.NewUserUseCase(userRepository)
	userHandler := handler.NewUserHandler(userUseCase)
	serverHTTP := api.NewServerHTTP(userHandler)
	return serverHTTP, nil
}
