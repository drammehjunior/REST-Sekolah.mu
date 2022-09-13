//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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
	gorpDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepository(gorpDB)
	userUseCase := usecase.NewUserUseCase(userRepository)
	userHandler := handler.NewUserHandler(userUseCase)
	serverHTTP := api.NewServerHTTP(userHandler)
	return serverHTTP, nil
}
