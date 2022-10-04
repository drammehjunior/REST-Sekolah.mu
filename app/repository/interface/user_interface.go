package _interface

import (
	"exampleclean.com/refactor/app/domain"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
)

type UserRepository interface {
	FindAll() ([]domain.Users, error)
	FindByID(id uint) *domain.Users
	FindByEmail(email string) (*domain.Users, error)
	Save(user rest_structs.RequestSignup) error
	Delete(user domain.Users) error
	UpdatePassword(user domain.Users) (int64, error)
}
