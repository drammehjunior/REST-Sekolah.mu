package mocks

import (
	"errors"
	"exampleclean.com/refactor/app/domain"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
	"fmt"
	"github.com/stretchr/testify/mock"
)

type UserUsecaseMock struct {
	Mock mock.Mock
}

func (usecase *UserUsecaseMock) FindByID(id uint) (*domain.Users, error) {
	args := usecase.Mock.Called(id)

	if args.Get(0) != nil {
		user := args.Get(0).(domain.Users)
		return &user, nil
	} else {
		return nil, errors.New("")
	}

}

func (usecase *UserUsecaseMock) FindByEmail(email string) (*domain.Users, error) {
	args := usecase.Mock.Called(email)

	fmt.Println("-----------")
	if args.Get(0) != nil {
		user := args.Get(0).(domain.Users)
		return &user, nil
	} else {
		return nil, errors.New("")
	}

}

func (usecase *UserUsecaseMock) Delete(id uint) error {
	args := usecase.Mock.Called(id)

	if args.Get(0) != nil {
		return errors.New("")

	} else {
		return nil
	}
}

func (usecase *UserUsecaseMock) UpdatePassword(user rest_structs.UpdatePassword) error {
	args := usecase.Mock.Called(user)
	if args.Get(0) != nil {
		return errors.New("")
	} else {
		return nil
	}
}

func (usecase *UserUsecaseMock) Login(user rest_structs.LoginBody) (*domain.Users, string, error) {
	args := usecase.Mock.Called(user)
	var userGet domain.Users
	var token string
	var err error

	if args.Get(0) != nil && args.Get(1) != nil {
		userGet = args.Get(0).(domain.Users)
		token = args.Get(1).(string)
		err = nil
		return &userGet, token, err
	}

	if args.Get(2) != nil {
		return nil, "", errors.New("fail to login user")
	}

	return nil, "", nil
}

func (usecase *UserUsecaseMock) Save(user rest_structs.RequestSignup) (userGet *rest_structs.SignupResponse, err error) {
	args := usecase.Mock.Called(user)

	if temp, con := args.Get(0).(rest_structs.SignupResponse); con {
		userGet = &temp
		//fmt.Println("Hello")
	} else {
		userGet = nil
	}

	if temp, con := args.Get(1).(func() error); con {
		err = temp()
	} else {
		err = args.Error(1)
	}

	return
}

func (usecase *UserUsecaseMock) FindAll() ([]domain.Users, error) {
	args := usecase.Mock.Called()

	if args.Get(0) != nil {
		return args.Get(0).([]domain.Users), nil
	}

	if args.Get(1).(error) != nil {
		return nil, args.Get(1).(error)
	}

	return nil, nil
}
