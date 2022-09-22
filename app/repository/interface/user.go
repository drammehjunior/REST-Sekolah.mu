package _interface

import (
	"exampleclean.com/refactor/app/domain"
)

type UserRepository interface {
	FindAll() ([]domain.Users, error)
	FindByID(id uint) (*domain.Users, error)
	FindByEmail(email string) (*domain.Users, error)
	Save(user domain.Users) (domain.Users, error)
	Delete(user domain.Users) error
	UpdatePassword(user domain.Users) (int64, error)
}
