package repository

import (
	"errors"
	"exampleclean.com/refactor/app/domain"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) FindAll() ([]domain.Users, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *UserRepositoryMock) FindByEmail(email string) (domain.Users, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *UserRepositoryMock) Save(user domain.Users) (domain.Users, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *UserRepositoryMock) Delete(user domain.Users) error {
	//TODO implement me
	panic("implement me")
}

func (repository *UserRepositoryMock) UpdatePassword(user domain.Users) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (repository *UserRepositoryMock) FindByID(id uint) (domain.Users, error) {
	arguments := repository.Mock.Called(id)
	if arguments.Get(0) == nil {
		return domain.Users{}, errors.New("error from FindById")
	} else {
		user := arguments.Get(0).(domain.Users)
		return user, nil
	}
}
