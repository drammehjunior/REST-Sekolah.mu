package usecase

import (
	"errors"
	"exampleclean.com/refactor/app/domain"
	interfaces "exampleclean.com/refactor/app/repository/interface"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
	_interface "exampleclean.com/refactor/app/usecase/interface"
	helper "exampleclean.com/refactor/app/utils"
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
	//c.userRepo.FindAll()
	return users, err
}

func (c *userUseCase) FindByID(id uint) (domain.Users, error) {
	user, err := c.userRepo.FindByID(id)
	return user, err
}

// implement login here
func (c *userUseCase) Save(user domain.Users) (domain.Users, error) {
	user, err := c.userRepo.Save(user)

	return user, err
}

func (c *userUseCase) Delete(user domain.Users) error {
	err := c.userRepo.Delete(user)
	return err
}

func (c *userUseCase) FindByEmail(email string) (*domain.Users, error) {
	user, err := c.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user account not found")
	} else {
		return user, nil
	}

}

func (c *userUseCase) UpdatePassword(user domain.Users) (int64, error) {
	row, err := c.userRepo.UpdatePassword(user)
	return row, err
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
