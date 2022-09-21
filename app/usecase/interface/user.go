package _interface

import (
	"exampleclean.com/refactor/app/domain"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
)

type UserUseCase interface {
	FindAll() ([]domain.Users, error)
	FindByID(id uint) (domain.Users, error)
	FindByEmail(email string) (*domain.Users, error)
	Save(user domain.Users) (domain.Users, error)
	Delete(user domain.Users) error
	UpdatePassword(user domain.Users) (int64, error)
	Login(user rest_structs.LoginBody) (*domain.Users, string, error)
}
