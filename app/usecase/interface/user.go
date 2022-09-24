package _interface

import (
	"exampleclean.com/refactor/app/domain"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
)

type UserUseCase interface {
	FindAll() ([]domain.Users, error)
	FindByID(id uint) (*domain.Users, error)
	FindByEmail(email string) (*domain.Users, error)
	Save(user rest_structs.RequestSignup) (*rest_structs.SignupResponse, error)
	Delete(id uint) error
	UpdatePassword(user rest_structs.UpdatePassword) error
	Login(user rest_structs.LoginBody) (*domain.Users, string, error)
}
