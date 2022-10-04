package repository

import (
	"errors"
	"exampleclean.com/refactor/app/domain"
	"exampleclean.com/refactor/app/repository/interface"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
	helper "exampleclean.com/refactor/app/utils"
	"github.com/jinzhu/gorm"
)

type userDatabaseGorm struct {
	DB *gorm.DB
}

func NewUserRepositoryGorm(DB *gorm.DB) _interface.UserRepository {
	return &userDatabaseGorm{DB}
}

func (c *userDatabaseGorm) FindAll() ([]domain.Users, error) {
	var users []domain.Users

	res := c.DB.Find(&users)
	//fmt.Println(users)

	if err := res.Error; err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("users are empty")
	} else {
		return users, nil
	}
}

func (c *userDatabaseGorm) FindByID(id uint) (user *domain.Users) {
	var users domain.Users

	temp := int(id)
	res := c.DB.Find(&users, "Id = ?", temp)

	if res.Error != nil {
		return nil
	}

	return &users

}

func (c *userDatabaseGorm) Save(user rest_structs.RequestSignup) error {
	parsedUser := domain.Users{
		Email:     user.Email,
		Password:  helper.HashThisSHA1(user.Password),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
	}
	res := c.DB.Create(&parsedUser)

	if err := res.Error; err != nil {
		return err
	} else {
		return nil
	}
}

func (c *userDatabaseGorm) Delete(user domain.Users) error {
	res := c.DB.Delete(&domain.Users{}, user.Id)
	if res.Error != nil {
		return errors.New("Internal error, cannot delete user")
	}
	return nil
}

func (c *userDatabaseGorm) FindByEmail(email string) (*domain.Users, error) {
	var user domain.Users
	res := c.DB.First(&user, "Email = ?", email)
	if res.Error != nil {
		return nil, errors.New("cannot find user with that email")
	}
	return &user, nil
}

func (c *userDatabaseGorm) UpdatePassword(user domain.Users) (int64, error) {
	user.Password = helper.HashThisSHA1(user.Password)

	res := c.DB.Model(&domain.Users{}).Where("Email = ?", user.Email).Update("Password", user.Password)
	if res.Error != nil {
		return 0, errors.New("failed to update user")
	}
	return res.RowsAffected, nil
}
