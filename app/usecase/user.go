package usecase

import (
	"errors"
	"exampleclean.com/refactor/app/domain"
	interfaces "exampleclean.com/refactor/app/repository/interface"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
	_interface "exampleclean.com/refactor/app/usecase/interface"
	helper "exampleclean.com/refactor/app/utils"
	"github.com/jinzhu/copier"
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
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *userUseCase) FindByID(id uint) (*domain.Users, error) {
	user := c.userRepo.FindByID(id)
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (c *userUseCase) Save(user rest_structs.RequestSignup) (*domain.Users, error) {
	//check if all the fields are filled
	var userParsed domain.Users

	if err := helper.IsEmailValid(user.Email); err != nil {
		return nil, err
	}

	if user.Firstname == "" || user.Lastname == "" {
		return nil, errors.New("first and last name are empty")
	}

	if _, err := c.FindByEmail(user.Email); err == nil {
		return nil, errors.New("user already exist. Please login")
	}

	hashedPassword, err := user.ValidateAndHash()
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	err = copier.Copy(&userParsed, &user)
	if err != nil {
		return nil, err
	}

	savedUser, err := c.userRepo.Save(userParsed)
	if err != nil {
		return nil, err
	}

	return &savedUser, nil
}

func (c *userUseCase) Delete(id uint) error {
	user := c.userRepo.FindByID(id)
	if user == nil {
		return errors.New("user not found")
	}

	if err := c.userRepo.Delete(*user); err != nil {
		return err
	}
	return nil
}

func (c *userUseCase) FindByEmail(email string) (*domain.Users, error) {
	if err := helper.IsEmailValid(email); err != nil {
		return nil, err
	}

	user, err := c.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user account not found")
	} else {
		return user, nil
	}

}

func (c *userUseCase) UpdatePassword(user rest_structs.UpdatePassword) error {
	//check if email and password are completed
	if err := helper.IsEmailValid(user.Email); err != nil {
		return err
	}

	if user.OldPassword == "" || user.NewPassword == "" || user.NewPasswordConfirm == "" {
		return errors.New("password fields cannot be empty")
	}

	hashedPassword, err := user.ValidateAndHash()
	if err != nil {
		return err
	}

	//check if the user account exist

	userNew, err := c.userRepo.FindByEmail(user.Email)
	if err != nil {
		return errors.New("user cannot be found")
	}

	//check if the previous passwords are matched
	if helper.IsPasswordMatched(userNew.Password, user.OldPassword) {
		return errors.New("passwords is incorrect")
	} else {
		userNew.Password = hashedPassword
	}

	//change password
	if _, err := c.userRepo.UpdatePassword(*userNew); err != nil {
		return err
	}

	return nil
}

func (c *userUseCase) Login(user rest_structs.LoginBody) (*domain.Users, string, error) {
	//check if both email and password are not empty
	if err := helper.IsLoginInputValid(user); err != nil {
		return nil, "", err
	}

	// check if user exist
	userNew, err := c.userRepo.FindByEmail(user.Email)
	if err != nil {
		return nil, "", errors.New("user not found")
	}

	// compare the passwords
	if helper.IsPasswordMatched(userNew.Password, user.Password) {
		return nil, "", errors.New("email or password is incorrect")
	}

	// if passwords are matched then make token and send
	token, _ := helper.SignInToken(*userNew)
	if token == "" {
		return nil, "", errors.New("error with token creation")
	}

	return userNew, token, nil

}
