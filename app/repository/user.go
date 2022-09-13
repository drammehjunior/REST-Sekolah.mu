package repository

import (
	"context"
	"exampleclean.com/refactor/app/domain"
	"exampleclean.com/refactor/app/repository/interface"
	"github.com/go-gorp/gorp"
)

type userDatabase struct {
	DB *gorp.DbMap
}

func NewUserRepository(DB *gorp.DbMap) _interface.UserRepository {
	return &userDatabase{DB}
}

func (c *userDatabase) FindAll(ctx context.Context) ([]domain.Users, error) {
	var users []domain.Users
	_, err := c.DB.Select(&users, "SELECT * FROM user")

	return users, err
}

func (c *userDatabase) FindByID(ctx context.Context, id uint) (domain.Users, error) {
	var user domain.Users
	err := c.DB.SelectOne(&user, "SELECT * FROM user WHERE Id=?", id)
	return user, err
}

func (c *userDatabase) Save(ctx context.Context, user domain.Users) (domain.Users, error) {
	err := c.DB.Insert(&user)

	return user, err
}

func (c *userDatabase) Delete(ctx context.Context, user domain.Users) error {
	_, err := c.DB.Delete(&user)
	return err
}

func (c *userDatabase) FindByEmail(ctx context.Context, email string) (domain.Users, error) {
	var user domain.Users
	err := c.DB.SelectOne(&user, "SELECT * FROM user WHERE Email=?", email)
	return user, err
}

func (c *userDatabase) UpdatePassword(ctx context.Context, user domain.Users) (int64, error) {
	row, err := c.DB.Update(&user)
	return row, err
}
