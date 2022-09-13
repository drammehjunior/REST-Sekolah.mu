package usecase

import (
	"context"
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

func (c *userUseCase) FindAll(ctx context.Context) ([]domain.Users, error) {
	users, err := c.userRepo.FindAll(ctx)
	c.userRepo.FindAll(ctx)
	return users, err
}

func (c *userUseCase) FindByID(ctx context.Context, id uint) (domain.Users, error) {
	user, err := c.userRepo.FindByID(ctx, id)
	return user, err
}

func (c *userUseCase) Save(ctx context.Context, user domain.Users) (domain.Users, error) {
	user, err := c.userRepo.Save(ctx, user)

	return user, err
}

func (c *userUseCase) Delete(ctx context.Context, user domain.Users) error {
	err := c.userRepo.Delete(ctx, user)
	return err
}

func (c *userUseCase) FindByEmail(ctx context.Context, email string) (domain.Users, error) {
	user, err := c.userRepo.FindByEmail(ctx, email)
	return user, err
}

func (c *userUseCase) UpdatePassword(ctx context.Context, user domain.Users) (int64, error) {
	row, err := c.userRepo.UpdatePassword(ctx, user)
	return row, err
}
