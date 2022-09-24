package repository

import (
	"errors"
	"exampleclean.com/refactor/app/domain"
	"exampleclean.com/refactor/app/repository/interface"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
	helper "exampleclean.com/refactor/app/utils"
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

func (c *userDatabase) Save(user rest_structs.RequestSignup) error {
	parsedUser := domain.Users{
		Email:     user.Email,
		Password:  helper.HashPassword(user.Password),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
	}
	err := c.DB.Insert(&parsedUser)
	if err != nil {
		return err
	} else {
		return nil
	}
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
	user.Password = helper.HashPassword(user.Password)
	row, err := c.DB.Update(&user)
	return row, err
}
