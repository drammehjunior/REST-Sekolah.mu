package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"exampleclean.com/refactor/app/domain"
	usecase2 "exampleclean.com/refactor/app/mocks"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
	helper "exampleclean.com/refactor/app/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetRouter() *gin.Engine {
	router := gin.Default()
	return router
}

var usecase = usecase2.UserUsecaseMock{Mock: mock.Mock{}}
var handler = UserHandler{userUseCase: &usecase}

func TestUserHandler_DeleteSuccess(t *testing.T) {

	type UserId struct {
		ID uint `json:"id"`
	}

	r := SetRouter()
	r.DELETE("/delete/:id", handler.Delete)

	userIdUser := UserId{ID: uint(4)}

	usecase.Mock.On("Delete", userIdUser.ID).Return(nil).Once()

	jsonData, _ := json.Marshal(userIdUser)
	req, _ := http.NewRequest("DELETE", fmt.Sprintf(`/delete/%d`, userIdUser.ID), bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestUserHandler_DeleteFailed(t *testing.T) {

	type RespondCall struct {
		Error string `json:"error"`
	}

	type UserId struct {
		ID uint `json:"id"`
	}

	r := SetRouter()
	r.DELETE("/delete/:id", handler.Delete)

	userIdUser := UserId{ID: uint(4)}

	t.Run("user_test_1", func(t *testing.T) {

		usecase.Mock.On("Delete", userIdUser.ID).Return(errors.New("")).Once()

		jsonData, _ := json.Marshal(userIdUser)
		req, _ := http.NewRequest("DELETE", fmt.Sprintf(`/delete/%s`, "dd"), bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		res, _ := io.ReadAll(w.Body)
		response := RespondCall{}
		json.Unmarshal(res, &response)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "cannot parse id", response.Error)
	})

	t.Run("user_test_2", func(t *testing.T) {
		usecase.Mock.On("Delete", userIdUser.ID).Return(errors.New("")).Once()

		jsonData, _ := json.Marshal(userIdUser)
		req, _ := http.NewRequest("DELETE", fmt.Sprintf(`/delete/%d`, userIdUser.ID), bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		res, _ := io.ReadAll(w.Body)
		response1 := RespondCall{}
		json.Unmarshal(res, &response1)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "failed to delete user", response1.Error)
	})

}

func TestUserHandler_FindByEmailSuccess(t *testing.T) {
	type RespondCall struct {
		Email     string `json:"Email"`
		Firstname string `json:"Firstname"`
		Lastname  string `json:"Lastname"`
	}

	var userEmail = "skeolahmu@gmail.com"

	user := domain.Users{
		Id:        uint(2),
		Email:     userEmail,
		Password:  "123454",
		Firstname: "sekolah",
		Lastname:  "sekolah",
	}

	r := SetRouter()
	r.GET("/user/email/:mail", handler.FindByEmail)

	usecase.Mock.On("FindByEmail", userEmail).Return(user, nil).Once()

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest("GET", fmt.Sprintf(`/user/email/%s`, userEmail), bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	respondBody, _ := ioutil.ReadAll(w.Body)

	//fmt.Println(string(respondBody))
	var respondCall RespondCall
	json.Unmarshal(respondBody, &respondCall)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, user.Email, respondCall.Email)
	assert.Equal(t, user.Firstname, respondCall.Firstname)
	assert.Equal(t, user.Lastname, respondCall.Lastname)
}

func TestUserHandler_FindByEmailFailed(t *testing.T) {
	type RespondCall struct {
		Error string `json:"error"`
	}

	var userEmail = "skeolahmu@gmail.com"

	user := domain.Users{
		Id:        uint(2),
		Email:     userEmail,
		Password:  "123454",
		Firstname: "sekolah",
		Lastname:  "sekolah",
	}

	r := SetRouter()
	r.GET("/user/email/:mail", handler.FindByEmail)

	t.Run("user_test_1", func(t *testing.T) {
		usecase.Mock.On("FindByEmail", userEmail).Return(nil, errors.New("")).Once()

		jsonData, _ := json.Marshal(user)
		req, _ := http.NewRequest("GET", fmt.Sprintf(`/user/email/%s`, userEmail), bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		respondBody, _ := ioutil.ReadAll(w.Body)

		//fmt.Println(string(respondBody))
		var respondCall RespondCall
		json.Unmarshal(respondBody, &respondCall)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("user_test_2", func(t *testing.T) {
		usecase.Mock.On("FindByEmail", userEmail).Return(nil, errors.New("")).Once()

		jsonData, _ := json.Marshal(user)
		req, _ := http.NewRequest("GET", fmt.Sprintf(`/user/email/%d`, 3), bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		respondBody, _ := ioutil.ReadAll(w.Body)

		//fmt.Println(string(respondBody))
		var respondCall RespondCall
		json.Unmarshal(respondBody, &respondCall)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "parameter is invalid", respondCall.Error)
	})

}

func TestUserHandler_FindByIDSuccess(t *testing.T) {
	var userId = uint(3)

	user := domain.Users{
		Id:        userId,
		Email:     "sekolahmu@gmail.com",
		Password:  "123454",
		Firstname: "sekolah",
		Lastname:  "sekolah",
	}

	r := SetRouter()
	r.GET("/users/:id", handler.FindByID)

	usecase.Mock.On("FindByID", userId).Return(user, nil).Once()

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest("GET", "/users/3", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	respondBody, _ := ioutil.ReadAll(w.Body)
	//fmt.Println(string(respondBody))

	var respond Response
	json.Unmarshal(respondBody, &respond)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestUserHandler_FindByID(t *testing.T) {
	var userId = uint(3)

	user := domain.Users{
		Id:        userId,
		Email:     "sekolahmu@gmail.com",
		Password:  "123454",
		Firstname: "sekolah",
		Lastname:  "sekolah",
	}

	r := SetRouter()
	r.GET("/users/:id", handler.FindByID)

	usecase.Mock.On("FindByID", userId).Return(nil, errors.New("")).Once()

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest("GET", "/users/3", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	respondBody, _ := ioutil.ReadAll(w.Body)
	//fmt.Println(string(respondBody))

	var respond Response
	json.Unmarshal(respondBody, &respond)

	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestUserHandler_LoginHandlerSuccess(t *testing.T) {

	type Response struct {
		Data struct {
			Email     string `json:"Email"`
			Password  string `json:"Password"`
			Firstname string `json:"Firstname"`
			Lastname  string `json:"Lastname"`
		} `json:"Data"`
		Token string `json:"Token"`
		Error bool   `json:"error"`
	}

	user := domain.Users{
		Id:        uint(8),
		Email:     "sekolahmu1@gmail.com",
		Password:  "12345",
		Firstname: "sekoa",
		Lastname:  "eeee",
	}

	userLogin := rest_structs.LoginBody{
		Email:    "sekolahmu1@gmail.com",
		Password: "12345",
	}

	r := SetRouter()
	r.POST("/login", handler.LoginHandler)

	usecase.Mock.On("Login", userLogin).Return(user, "1234", nil).Once()

	jsonData, _ := json.Marshal(userLogin)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	respondBody, _ := ioutil.ReadAll(w.Body)
	//fmt.Println(string(respondBody))

	var respond Response
	json.Unmarshal(respondBody, &respond)

	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Equal(t, respond.Error, false)
	assert.Equal(t, respond.Data.Email, userLogin.Email)
}

func TestUserHandler_LoginHandlerFailed(t *testing.T) {
	r := SetRouter()
	r.POST("/login", handler.LoginHandler)

	type Response struct {
		Error string `json:"error"`
	}

	t.Run("user_test_1", func(t *testing.T) {
		userLogin := rest_structs.LoginBody{
			Email:    "sekolahmu1@gmail.com",
			Password: "12345",
		}

		usecase.Mock.On("Login", userLogin).Return(nil, nil, errors.New("")).Once()

		jsonData, _ := json.Marshal(userLogin)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var respond Response
		respondBody, _ := ioutil.ReadAll(w.Body)

		json.Unmarshal(respondBody, &respond)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, respond.Error, "fail to login user")

	})

	t.Run("user_test_2", func(t *testing.T) {
		userLogin1 := rest_structs.LoginBody{
			Email:    "sekolahmu1@gmail.com",
			Password: "",
		}

		usecase.Mock.On("Login", userLogin1).Return(nil, nil, errors.New("")).Once()

		jsonData, _ := json.Marshal(userLogin1)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var respond Response
		respondBody, _ := ioutil.ReadAll(w.Body)

		json.Unmarshal(respondBody, &respond)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, respond.Error, "email or password cannot be empty")

	})

}

func TestUserHandler_UpdatePasswordFailed(t *testing.T) {
	r := SetRouter()
	r.PUT("/updatePassword", handler.UpdatePassword)

	user := rest_structs.UpdatePassword{
		Email:              "drammmeh@gmail.com",
		OldPassword:        "1234",
		NewPassword:        "12345",
		NewPasswordConfirm: "12345",
	}

	t.Run("user_test_1", func(t *testing.T) {
		usecase.Mock.On("UpdatePassword", user).Return(errors.New("")).Once()

		jsonData, _ := json.Marshal(user)
		req, _ := http.NewRequest("PUT", "/updatePassword", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("user_test_2", func(t *testing.T) {
		usecase.Mock.On("UpdatePassword", user).Return().Once()

		jsonData, _ := json.Marshal(user)
		req, _ := http.NewRequest("PUT", "/updatePassword", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}

func TestUserHandler_UpdatePasswordSuccess(t *testing.T) {
	user := rest_structs.UpdatePassword{
		Email:              "drammeh@gmail.com",
		OldPassword:        "1234",
		NewPassword:        "12345",
		NewPasswordConfirm: "12345",
	}
	r := SetRouter()
	r.PUT("/updatePassword", handler.UpdatePassword)

	usecase.Mock.On("UpdatePassword", user).Return(nil).Once()

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/updatePassword", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestUserHandler_FindAllFailed(t *testing.T) {
	r := SetRouter()
	r.GET("/users", handler.FindAll)

	usecase.Mock.On("FindAll").Return(nil, errors.New("temp")).Once()
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestUserHandler_FindAllSuccess(t *testing.T) {
	r := SetRouter()
	r.GET("/users", handler.FindAll)

	type Response struct {
		Data []struct {
			ID        int    `json:"Id"`
			Email     string `json:"Email"`
			Password  string `json:"Password"`
			Firstname string `json:"Firstname"`
			Lastname  string `json:"Lastname"`
		} `json:"data"`
		Result int `json:"result"`
	}

	output := []domain.Users{
		{
			Id:        uint(12),
			Email:     "sekolahmu3@gmail.com",
			Password:  helper.HashPassword("1234"),
			Firstname: "sekolah",
			Lastname:  "mu",
		},
		{
			Id:        uint(13),
			Email:     "sekolahmu4@gmail.com",
			Password:  helper.HashPassword("1234"),
			Firstname: "sekolah",
			Lastname:  "mu",
		},
	}

	usecase.Mock.On("FindAll").Return(output, nil)

	jsonValue, _ := json.Marshal(output)
	req, _ := http.NewRequest("GET", "/users", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res, _ := ioutil.ReadAll(w.Body)

	var respond Response

	//fmt.Println(string(res))

	json.Unmarshal(res, &respond)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, len(respond.Data), len(output))

	for index, data := range respond.Data {
		//fmt.Println(data)
		//fmt.Println(output[index])

		t.Run(fmt.Sprintf(`test user %d`, index), func(t *testing.T) {
			assert.Equal(t, data.Email, output[index].Email)
			assert.Equal(t, data.Password, output[index].Password)
			assert.Equal(t, data.Lastname, output[index].Lastname)
			assert.Equal(t, data.Firstname, output[index].Firstname)
		})
	}
}

func TestUserHandler_SaveSignupFailed(t *testing.T) {

	type RespondCall struct {
		Error string `json:"error"`
	}

	r := SetRouter()
	r.POST("/signup", handler.SaveSignup)

	input := rest_structs.RequestSignup{
		Email:           "sekolahmu1@gmail.com",
		Password:        "12345",
		PasswordConfirm: "12345",
		Firstname:       "sekolah",
		Lastname:        "mu",
	}

	t.Run("fail case 1", func(t *testing.T) {
		usecase.Mock.On("Save", input).Return(nil, errors.New(""))

		jsonValue, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		//
		//fmt.Println(w.Body)
	})

	t.Run("fail case 1", func(t *testing.T) {

		input1 := rest_structs.RequestSignup{
			Email:    "sekolahmu1@gmail.com",
			Password: "12345",
		}

		usecase.Mock.On("Save", input1).Return(nil, errors.New(""))

		jsonValue, _ := json.Marshal(input1)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res, _ := ioutil.ReadAll(w.Body)

		fmt.Println(string(res))

		//json.Unmarshal(res, &response)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		//assert.Equal(t, "problem with the body", response.Error)
		//fmt.Println(w.Body)
	})

}

func TestUserHandler_SaveSignupSuccess(t *testing.T) {
	input := rest_structs.RequestSignup{
		Email:           "sekolahmu2@gmail.com",
		Password:        "12345",
		PasswordConfirm: "12345",
		Firstname:       "sekolah",
		Lastname:        "mu",
	}

	respone := rest_structs.SignupResponse{
		Email:     input.Email,
		Firstname: input.Firstname,
		Lastname:  input.Lastname,
	}

	mockResponse := `{
    "data": {
        "Email": "%s",
        "Firstname": "%s",
        "Lastname": "%s"
    	}
	}`

	mockResponse = fmt.Sprintf(mockResponse, input.Email, input.Firstname, input.Lastname)

	r := SetRouter()
	r.POST("/signup", handler.SaveSignup)

	usecase.Mock.On("Save", input).Return(respone, nil)

	jsonRequest, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonRequest))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, http.StatusCreated, w.Code)
	fmt.Println(string(res))
}
