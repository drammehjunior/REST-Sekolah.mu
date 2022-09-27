package handler

import (
	structs "exampleclean.com/refactor/app/rest-structs"
	services "exampleclean.com/refactor/app/usecase/interface"
	"exampleclean.com/refactor/app/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
)

var jwtKey = []byte("my_super_secret_key")

type UserHandler struct {
	userUseCase services.UserUseCase
}

type Response struct {
	Id        int64  `copier:"must"`
	Email     string `copier:"must"`
	Password  string `copier:"must"`
	Firstname string `copier:"must"`
	Lastname  string `copier:"must"`
}

type FindByEmailResponse struct {
	Email     string `copier:"must"`
	Firstname string `copier:"must"`
	Lastname  string `copier:"must"`
}

type ResponseLogin struct {
	Email     string `copier:"must"`
	Password  string `copier:"must"`
	Firstname string `copier:"must"`
	Lastname  string `copier:"must"`
}

func NewUserHandler(usercase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usercase,
	}
}

func (cr *UserHandler) FindAll(c *gin.Context) {
	users, err := cr.userUseCase.FindAll()

	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		response := []Response{}
		copier.Copy(&response, &users)
		c.JSON(http.StatusOK, gin.H{
			"data":   response,
			"result": len(response),
		})

	}
}

func (cr *UserHandler) FindByID(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot parse id",
		})
		return
	}
	user, err := cr.userUseCase.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
	} else {
		response := Response{}
		copier.Copy(&response, &user)
		c.JSON(http.StatusOK, response)
	}
}

func (cr *UserHandler) SaveSignup(c *gin.Context) {
	var userSignup structs.RequestSignup

	//get the body
	c.BindJSON(&userSignup)

	user, err := cr.userUseCase.Save(userSignup)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": user})
	return
}

func (cr *UserHandler) Delete(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot parse id",
		})
		return
	}

	if err := cr.userUseCase.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete user",
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "user deleted successfully",
	})
	return
}

func (cr *UserHandler) FindByEmail(c *gin.Context) {
	paramsEmail := c.Param("mail")

	if err := utils.IsEmailValid(paramsEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "parameter is invalid",
		})
		return
	}

	user, err := cr.userUseCase.FindByEmail(paramsEmail)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "parameter is invalid",
		})
		return
	} else {
		res := FindByEmailResponse{}
		copier.Copy(&res, &user)
		c.JSON(http.StatusOK, res)
		return
	}
}

func (cr *UserHandler) LoginHandler(c *gin.Context) {

	var requestBody structs.LoginBody

	c.BindJSON(&requestBody)

	if requestBody.Email == "" || requestBody.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or password cannot be empty"})
		return
	}

	user, token, err := cr.userUseCase.Login(requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := ResponseLogin{}
	copier.Copy(&response, &user)

	c.JSON(http.StatusAccepted, gin.H{
		"error": false,
		"Data":  response,
		"Token": token,
	})
	return
}

// UpdatePassword should be implemented after the use of token
func (cr *UserHandler) UpdatePassword(c *gin.Context) {
	var requestBody structs.UpdatePassword

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Internal Error"})
		return
	}

	err := cr.userUseCase.UpdatePassword(requestBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "password changed",
	})
	return
}
