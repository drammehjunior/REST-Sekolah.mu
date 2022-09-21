package repository

import (
	"errors"
	"exampleclean.com/refactor/app/domain"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) Delete(user domain.Users) error {
	//TODO implement me
	panic("implement me")
}

func (repository *UserRepositoryMock) UpdatePassword(user domain.Users) (int64, error) {
	arguments := repository.Mock.Called("UpdatePassword")

	if arguments.Get(0) == nil {
		return 0, errors.New("error found in UpdatePassword")
	} else {
		return 1, nil
	}
}

func (repository *UserRepositoryMock) FindByEmail(email string) (domain.Users, error) {
	arguments := repository.Mock.Called(email)

	if arguments.Get(0) == nil {
		return domain.Users{}, errors.New("error found in FindByEmail")
	} else {
		return arguments.Get(0).(domain.Users), nil
	}
}

func (repository *UserRepositoryMock) FindAll() ([]domain.Users, error) {

	arguments := repository.Mock.Called()

	if arguments.Get(0) == nil {
		return []domain.Users{}, errors.New("error found in FindAll")
	} else {
		user := arguments.Get(0).([]domain.Users)
		return user, nil
	}
}

func (repository *UserRepositoryMock) FindByID(id uint) (domain.Users, error) {
	arguments := repository.Mock.Called(id)

	//fmt.Println(arguments.Is(domain.Users{}, domain.Users{}))
	//fmt.Printf("type: %#v", arguments.Get(0).(domain.Users))
	if arguments.Get(0) == nil {
		return domain.Users{}, errors.New("error from FindById")
	} else {
		user := arguments.Get(0).(domain.Users)
		return user, nil
	}
}

func (repository *UserRepositoryMock) Save(user domain.Users) (domain.Users, error) {
	arguments := repository.Mock.Called(user)

	//fmt.Println(arguments.Is(domain.Users{}, domain.Users{}))
	//fmt.Printf("type: %#v", arguments.Get(0).(domain.Users))
	if arguments.Get(0) == nil {
		return domain.Users{}, errors.New("error from FindById")
	} else {
		user := arguments.Get(0).(domain.Users)
		return user, nil
	}
}
