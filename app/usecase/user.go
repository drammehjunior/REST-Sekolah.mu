package usecase

import (
	"exampleclean.com/refactor/app/domain"
	interfaces "exampleclean.com/refactor/app/repository/interface"
	_interface "exampleclean.com/refactor/app/usecase/interface"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) _interface.UserUseCase {

	return &userUseCase{
		userRepo: repo,
	}
}

func (c *userUseCase) FindAll() ([]domain.Users, error) {
	users, err := c.userRepo.FindAll()
	c.userRepo.FindAll()
	return users, err
}

func (c *userUseCase) FindByID(id uint) (domain.Users, error) {
	user, err := c.userRepo.FindByID(id)
	return user, err
}

func (c *userUseCase) Save(user domain.Users) (domain.Users, error) {
	user, err := c.userRepo.Save(user)

	return user, err
}

func (c *userUseCase) Delete(user domain.Users) error {
	err := c.userRepo.Delete(user)
	return err
}

func (c *userUseCase) FindByEmail(email string) (domain.Users, error) {
	user, err := c.userRepo.FindByEmail(email)
	return user, err
}

func (c *userUseCase) UpdatePassword(user domain.Users) (int64, error) {
	row, err := c.userRepo.UpdatePassword(user)
	return row, err
}
