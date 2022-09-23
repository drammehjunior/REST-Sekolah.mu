package repository

import (
	"errors"
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

func (c *userDatabase) FindAll() ([]domain.Users, error) {
	var users []domain.Users
	_, err := c.DB.Select(&users, "SELECT * FROM user")

	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("users not found")
	} else {
		return users, nil
	}
}

func (c *userDatabase) FindByID(id uint) *domain.Users {
	var user domain.Users
	err := c.DB.SelectOne(&user, "SELECT * FROM user WHERE Id=?", id)
	if err != nil {
		return nil
	}
	return &user
}

func (c *userDatabase) Save(user domain.Users) (domain.Users, error) {
	err := c.DB.Insert(&user)

	return user, err
}

func (c *userDatabase) Delete(user domain.Users) error {
	_, err := c.DB.Delete(&user)
	if err != nil {
		return err
	}
	return nil
}

func (c *userDatabase) FindByEmail(email string) (*domain.Users, error) {
	var user domain.Users
	err := c.DB.SelectOne(&user, "SELECT * FROM user WHERE Email=?", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *userDatabase) UpdatePassword(user domain.Users) (int64, error) {
	row, err := c.DB.Update(&user)
	return row, err
}
