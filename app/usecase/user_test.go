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

	userRepository1.Mock.On("FindByID", uint(20)).Return(nil, nil).Once()
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

	userRepository1.Mock.On("FindByID", uint(2)).Return(user).Once()
	res, err := userUsecase.FindByID(2)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res)
	assert.Equal(t, user.Email, res.Email)
	userRepository1.Mock.AssertCalled(t, "FindByID", mock.AnythingOfType("uint"))

}

func TestUserUseCase_FindAllSuccess(t *testing.T) {
	users := []domain.Users{
		{
			Email:     "sekolahmu.mu@gmail.com",
			Password:  "1234556",
			Firstname: "med",
			Lastname:  "Dram",
		},
		{
			Email:     "sekolahmu22.mu@gmail.com",
			Password:  "1234556",
			Firstname: "med",
			Lastname:  "Dram",
		},
		{
			Email:     "sekolahmu33.mu@gmail.com",
			Password:  "1234556",
			Firstname: "med",
			Lastname:  "Dram",
		},
	}

	userRepository1.Mock.On("FindAll").Return(users, nil).Once()
	res, err := userUsecase.FindAll()

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, cap(res), cap(users))
	//fmt.Println(users)
}

func TestUserUseCase_FindAllFailed(t *testing.T) {
	userRepository1.Mock.On("FindAll").Return(nil, nil).Once()
	res, err := userUsecase.FindAll()

	assert.NotNil(t, err)
	assert.Empty(t, res)
}

func TestUserUseCase_FindByEmailSuccess(t *testing.T) {
	user := domain.Users{
		Id:        1,
		Email:     "sekolahmu1.mu@gmail.com",
		Password:  "1234556",
		Firstname: "med",
		Lastname:  "Dram",
	}

	email := "sekolahmu1@gmail.com"
	userRepository1.Mock.On("FindByEmail", email).Return(user, nil).Once()

	res, err := userUsecase.FindByEmail(email)
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.Equal(t, res.Email, user.Email)
}
