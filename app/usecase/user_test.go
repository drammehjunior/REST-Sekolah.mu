package usecase

import (
	"errors"
	"exampleclean.com/refactor/app/domain"
	"exampleclean.com/refactor/app/mocks"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
	helper "exampleclean.com/refactor/app/utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var userRepository1 = mocks.UserRepositoryMock{Mock: mock.Mock{}}
var userUsecase = userUseCase{&userRepository1}

func TestUserUseCase_FindByIDFailed(t *testing.T) {

	userRepository1.Mock.On("FindByID", uint(2000)).Return(nil).Once()
	user, err := userUsecase.FindByID(2000)
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Equal(t, err, errors.New("user not found"))

	//assert.Equal(t, "", user.Email)
}

func TestUserUseCase_FindByIDSuccess(t *testing.T) {

	userId := uint(23)
	user := domain.Users{
		Id:        userId,
		Email:     "sekolahmu.mu@gmail.com",
		Password:  "1234556",
		Firstname: "med",
		Lastname:  "Dram",
	}

	userRepository1.Mock.On("FindByID", userId).Return(user).Once()
	res, err := userUsecase.FindByID(userId)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.Id, userId)
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
}

func TestUserUseCase_FindAllFailed(t *testing.T) {
	userRepository1.Mock.On("FindAll").Return(nil).Once()
	user, err := userUsecase.FindAll()

	expectedErrorMessage := errors.New("users not found")

	assert.NotNil(t, err)
	assert.Empty(t, user)
	assert.Equal(t, err, expectedErrorMessage)
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

func TestUserUseCase_FindByEmailFailed(t *testing.T) {
	cases := []struct {
		name          string
		email         string
		expectedError error
	}{
		{
			name:          "user_test_1",
			email:         "mameddram.com",
			expectedError: errors.New("email is not valid"),
		},
		{
			name:          "user_test_2",
			email:         "",
			expectedError: errors.New("email cannot be empty"),
		},
	}

	for _, data := range cases {
		t.Run(data.name, func(t *testing.T) {
			userRepository1.Mock.On("FindByEmail", data.email).Return(nil).Once()
			user, err := userUsecase.FindByEmail(data.email)
			assert.Equal(t, data.expectedError, err)
			assert.Nil(t, user)
		})
	}

	userEmail := "sekolahmu1@gmail.com"
	expectedErr := errors.New("user account not found")

	t.Run("user_test_3", func(t *testing.T) {
		userRepository1.Mock.On("FindByEmail", userEmail).Return(nil).Once()
		user, err := userUsecase.FindByEmail(userEmail)
		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Equal(t, err, expectedErr)

	})
}

func TestUserUseCase_LoginSuccess(t *testing.T) {
	testCases := []struct {
		name    string
		request *rest_structs.LoginBody
	}{
		{
			name: "user_test_1",
			request: &rest_structs.LoginBody{
				Email:    "sekolahmu1@gmail.com",
				Password: "1234",
			},
		},
		{
			name: "user_test_2",
			request: &rest_structs.LoginBody{
				Email:    "sekolahmu2@gmail.com",
				Password: "1234",
			},
		},
	}

	for index, data := range testCases {
		testData := domain.Users{
			Id:        uint(index),
			Email:     data.request.Email,
			Password:  helper.HashPassword(data.request.Password),
			Firstname: "Sekolah",
			Lastname:  "Mu",
		}
		userRepository1.Mock.On("FindByEmail", testData.Email).Return(testData).Once()
		t.Run(data.name, func(t *testing.T) {
			user, token, err := userUsecase.Login(*data.request)
			assert.NotNil(t, user)
			assert.NotEmpty(t, token)
			assert.Nil(t, err)
		})

	}
}

func TestUserUseCase_LoginFailed(t *testing.T) {
	testCases := []struct {
		name          string
		request       *rest_structs.LoginBody
		expectedError error
	}{
		{
			name: "user_test_1",
			request: &rest_structs.LoginBody{
				Email:    "drammeh.com",
				Password: "1234",
			},
			expectedError: errors.New("email is not valid"),
		},
		{
			name: "user_test_2",
			request: &rest_structs.LoginBody{
				Email:    "sekolahmu1@gmail.com",
				Password: "",
			},
			expectedError: errors.New("email or password cannot be empty"),
		},
		{
			name: "user_test_3",
			request: &rest_structs.LoginBody{
				Email:    "",
				Password: "2344",
			},
			expectedError: errors.New("email or password cannot be empty"),
		},
		{
			name: "user_test_4",
			request: &rest_structs.LoginBody{
				Email:    "sekolahmu100@gmail.com",
				Password: "34455",
			},
			expectedError: errors.New("user not found"),
		},
	}

	for _, data := range testCases {
		t.Run(data.name, func(t *testing.T) {
			userRepository1.Mock.On("FindByEmail", data.request.Email).Return(nil, nil).Once()

			user, token, err := userUsecase.Login(*data.request)

			assert.NotNil(t, err)
			assert.Equal(t, err, data.expectedError)
			assert.Empty(t, token)
			assert.Nil(t, user)
		})
	}

	userDummyReturn := &domain.Users{
		Id:        3,
		Email:     "sekolahmutest@gmail.com",
		Password:  "1234",
		Firstname: "sekolah",
		Lastname:  "mu",
	}

	userInput := rest_structs.LoginBody{
		Email:    "sekolahmutest@gmail.com",
		Password: "12345",
	}

	expectedError := errors.New("email or password is incorrect")

	t.Run("user_test_5", func(t *testing.T) {
		userRepository1.Mock.On("FindByEmail", userInput.Email).Return(*userDummyReturn, nil).Once()
		user, token, err := userUsecase.Login(userInput)

		assert.Equal(t, err, expectedError)
		assert.Nil(t, user)
		assert.Empty(t, token)
	})

}

func TestUserUseCase_SaveSuccess(t *testing.T) {
	newUser := rest_structs.RequestSignup{
		Email:           "sekolahmu4@gmail.com",
		Password:        "1234",
		PasswordConfirm: "1234",
		Firstname:       "sekolah",
		Lastname:        "mu",
	}

	userRepository1.Mock.On("FindByEmail", newUser.Email).Return(nil).Once()
	userRepository1.Mock.On("Save", newUser).Return(nil).Once()
	user, err := userUsecase.Save(newUser)
	assert.Nil(t, err)
	assert.NotEmpty(t, user)
	assert.Equal(t, user.Email, newUser.Email)
}

func TestUserUseCase_SaveFailed(t *testing.T) {
	testCases := []struct {
		name          string
		request       *rest_structs.RequestSignup
		expectedError error
	}{
		{
			name: "user_test_1",
			request: &rest_structs.RequestSignup{
				Email:           "",
				Password:        "1234",
				PasswordConfirm: "1234",
				Firstname:       "mamed",
				Lastname:        "dram",
			},
			expectedError: errors.New("email cannot be empty"),
		},
		{
			name: "user_test_2",
			request: &rest_structs.RequestSignup{
				Email:           "sekolahmu.com",
				Password:        "1234",
				PasswordConfirm: "1234",
				Firstname:       "mamed",
				Lastname:        "dram",
			},
			expectedError: errors.New("email is not valid"),
		},
		{
			name: "user_test_3",
			request: &rest_structs.RequestSignup{
				Email:           "sekolahmu@gmail.com",
				Password:        "1234",
				PasswordConfirm: "1234",
				Firstname:       "",
				Lastname:        "dram",
			},
			expectedError: errors.New("first and last name are empty"),
		},
		{
			name: "user_test_4",
			request: &rest_structs.RequestSignup{
				Email:           "sekolahmu@gmail.com",
				Password:        "1234",
				PasswordConfirm: "1234",
				Firstname:       "mamed",
				Lastname:        "",
			},
			expectedError: errors.New("first and last name are empty"),
		},
		{
			name: "user_test_5",
			request: &rest_structs.RequestSignup{
				Email:           "sekolahmu@gmail.com",
				Password:        "1234",
				PasswordConfirm: "",
				Firstname:       "mamed",
				Lastname:        "dram",
			},
			expectedError: errors.New("passwords cannot be empty"),
		},
		{
			name: "user_test_6",
			request: &rest_structs.RequestSignup{
				Email:           "sekolahmu@gmail.com",
				Password:        "1234",
				PasswordConfirm: "12345",
				Firstname:       "mamed",
				Lastname:        "dram",
			},
			expectedError: errors.New("password do not match"),
		},
		{
			name: "user_test_6",
			request: &rest_structs.RequestSignup{
				Email:           "sekolahmu@gmail.com",
				Password:        "1234",
				PasswordConfirm: "1234",
				Firstname:       "mamed",
				Lastname:        "dram",
			},
			expectedError: errors.New("user already exist. Please login"),
		},
	}

	for _, data := range testCases {
		userTemplate := domain.Users{
			Email:     data.request.Email,
			Password:  data.request.Password,
			Firstname: data.request.Firstname,
			Lastname:  data.request.Lastname,
		}

		userSignup := rest_structs.RequestSignup{
			Email:           data.request.Email,
			Password:        data.request.Password,
			PasswordConfirm: data.request.PasswordConfirm,
			Firstname:       data.request.Firstname,
			Lastname:        data.request.Lastname,
		}

		t.Run(data.name, func(t *testing.T) {
			userRepository1.Mock.On("FindByEmail", data.request.Email).Return(userTemplate).Once()
			userRepository1.Mock.On("Save", userSignup).Return(nil).Once()

			user, err := userUsecase.Save(userSignup)
			assert.Equal(t, err, data.expectedError)
			assert.Nil(t, user)
		})

	}
}

func TestUserUseCase_Delete(t *testing.T) {
	userId := uint(8)
	user := domain.Users{
		Id:        userId,
		Email:     "sekolahmu@gmail.com",
		Password:  "12345",
		Firstname: "sekolah",
		Lastname:  "mu",
	}

	userRepository1.Mock.On("FindByID", userId).Return(user).Once()
	userRepository1.Mock.On("Delete", user).Return(nil).Once()
	err := userUsecase.Delete(userId)
	assert.Nil(t, err)
}

func TestUserUseCase_DeleteFailed(t *testing.T) {
	userId := uint(8)
	user := domain.Users{
		Id:        userId,
		Email:     "sekolahmu@gmail.com",
		Password:  "12345",
		Firstname: "sekolah",
		Lastname:  "mu",
	}

	t.Run("user_test_1", func(t *testing.T) {
		userRepository1.Mock.On("FindByID", userId).Return(nil).Once()
		userRepository1.Mock.On("Delete", user).Return(nil).Once()

		err := userUsecase.Delete(userId)
		assert.NotNil(t, err)
		fmt.Println(err)
		assert.Equal(t, err, errors.New("user not found"))
	})

	t.Run("user_test_2", func(t *testing.T) {

		user1 := domain.Users{
			Id:        uint(6),
			Email:     "sekolahmu1@gmail.com",
			Password:  "123451",
			Firstname: "sekolah",
			Lastname:  "mu",
		}

		userRepository1.Mock.On("FindByID", userId).Return(user1).Once()
		userRepository1.Mock.On("Delete", user1).Return(errors.New("")).Once()

		err := userUsecase.Delete(userId)
		assert.NotNil(t, err)
		fmt.Println(err)
		assert.Equal(t, err, errors.New("fail to delete user"))
	})

}

func TestUserUseCase_UpdatePasswordSuccess(t *testing.T) {
	userTest := rest_structs.UpdatePassword{
		Email:              "sekolahmue@gmail.com",
		OldPassword:        "1234",
		NewPassword:        "12345",
		NewPasswordConfirm: "12345",
	}

	userReturn := domain.Users{
		Id:        23,
		Email:     "sekolahmue@gmail.com",
		Password:  helper.HashPassword("1234"),
		Firstname: "sekolah",
		Lastname:  "mu",
	}
	userReturnChanged := domain.Users{
		Id:        23,
		Email:     "sekolahmue@gmail.com",
		Password:  userTest.NewPassword,
		Firstname: "sekolah",
		Lastname:  "mu",
	}

	t.Run("user_test_1", func(t *testing.T) {
		userRepository1.Mock.On("FindByEmail", userTest.Email).Return(userReturn).Once()
		userRepository1.Mock.On("UpdatePassword", userReturnChanged).Return(1).Once()

		err := userUsecase.UpdatePassword(userTest)
		assert.Nil(t, err)
	})
}

func TestUserUseCase_UpdatePasswordFailed(t *testing.T) {
	userTest := []struct {
		name          string
		request       *rest_structs.UpdatePassword
		expectedError error
	}{
		{
			name: "user_test_1",
			request: &rest_structs.UpdatePassword{
				Email:              "sekolahmu.com",
				OldPassword:        "1234",
				NewPassword:        "12345",
				NewPasswordConfirm: "12345",
			},
			expectedError: errors.New("email is not valid"),
		},
		{
			name: "user_test_2",
			request: &rest_structs.UpdatePassword{
				Email:              "sekolahmu@gmail.com",
				OldPassword:        "1234",
				NewPassword:        "12345",
				NewPasswordConfirm: "",
			},
			expectedError: errors.New("passwords cannot be empty"),
		},
		{
			name: "user_test_3",
			request: &rest_structs.UpdatePassword{
				Email:              "sekolahmu@gmail.com",
				OldPassword:        "",
				NewPassword:        "12345",
				NewPasswordConfirm: "12345",
			},
			expectedError: errors.New("password cannot be empty"),
		},
		{
			name: "user_test_3",
			request: &rest_structs.UpdatePassword{
				Email:              "sekolahmu@gmail.com",
				OldPassword:        "1222",
				NewPassword:        "12345",
				NewPasswordConfirm: "123456",
			},
			expectedError: errors.New("password do not match"),
		},
		{
			name: "user_test_4",
			request: &rest_structs.UpdatePassword{
				Email:              "",
				OldPassword:        "1234",
				NewPassword:        "12345",
				NewPasswordConfirm: "12345",
			},
			expectedError: errors.New("email cannot be empty"),
		},
	}

	for _, data := range userTest {
		userRepository1.Mock.On("FindByEmail", data.request.Email).Return(nil).Once()
		userRepository1.Mock.On("UpdatePassword", data.request).Return(nil).Once()
		err := userUsecase.UpdatePassword(*data.request)
		assert.NotNil(t, err)
		assert.Equal(t, err, data.expectedError)
	}
}
