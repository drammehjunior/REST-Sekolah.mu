package usecase

import (
	"exampleclean.com/refactor/app/domain"
	"exampleclean.com/refactor/app/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var userRepository1 = repository.UserRepositoryMock{Mock: mock.Mock{}}
var userUsecase = userUseCase{&userRepository1}

func TestUserUseCase_FindByIDFailed(t *testing.T) {

	userRepository1.Mock.On("FindByID", uint(20)).Return(nil, nil)

	user, err := userUsecase.FindByID(20)
	assert.NotNil(t, err)
	assert.Equal(t, "", user.Email)
	assert.Equal(t, "", user.Password)
	//assert.Equal(t, "", user.Email)
}

func TestUserUseCase_FindByIDSuccess(t *testing.T) {

	user := domain.Users{
		Email:     "sekolahmu.mu@gmail.com",
		Password:  "1234556",
		Firstname: "med",
		Lastname:  "Dram",
	}

	userRepository1.Mock.On("FindByID", uint(2)).Return(user)
	res, err := userUsecase.FindByID(2)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, user.Email, res.Email)
}
