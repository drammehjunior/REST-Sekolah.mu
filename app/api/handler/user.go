package handler

import (
	"exampleclean.com/refactor/app/domain"
	structs "exampleclean.com/refactor/app/rest-structs"
	services "exampleclean.com/refactor/app/usecase/interface"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
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

func NewUserHandler(usercase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usercase,
	}
}

func (cr *UserHandler) FindAll(c *gin.Context) {
	users, err := cr.userUseCase.FindAll(c.Request.Context())

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := []Response{}
		copier.Copy(&response, &users)
		c.JSON(http.StatusOK, response)

	}
}

func (cr *UserHandler) FindByID(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Cannot parse id",
		})
		return
	}
	user, err := cr.userUseCase.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := Response{}
		copier.Copy(&response, &user)
		c.JSON(http.StatusOK, response)
	}
}

func (cr *UserHandler) Save(c *gin.Context) {
	var user domain.Users

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}
	user, err := cr.userUseCase.Save(c.Request.Context(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := Response{}
		copier.Copy(&response, &user)
	}
}

func (cr *UserHandler) SaveSignup(c *gin.Context) {
	var user domain.Users
	var userSignup structs.RequestSignup

	//get the body
	if err := c.BindJSON(&userSignup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}
	//matches the passwords and return error if not match
	hashedPassword, errmsg := userSignup.ValidateAndHash()
	if errmsg != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "password does not match",
		})
		return
	}

	copier.Copy(&user, &userSignup)
	user.Password = hashedPassword

	user, err := cr.userUseCase.Save(c.Request.Context(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		response := Response{}
		copier.Copy(&response, &user)
		c.JSON(http.StatusOK, response)
	}
}

func (cr *UserHandler) Delete(c *gin.Context) {
	paramsId := c.Param("id")
	id, err := strconv.Atoi(paramsId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Cannot parse id",
		})
		return
	}

	ctx := c.Request.Context()
	user, err := cr.userUseCase.FindByID(ctx, uint(id))

	if user == (domain.Users{}) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User is cannot be found",
		})
		return
	}
	cr.userUseCase.Delete(ctx, user)
	c.JSON(http.StatusOK, gin.H{"message": "User is deleted succesfully"})
}

func (cr *UserHandler) FindByEmail(c *gin.Context) {
	paramsEmail := c.Param("mail")
	if paramsEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email cannot be empty",
		})
		return
	}
	user, err := cr.userUseCase.FindByEmail(c.Request.Context(), paramsEmail)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "email cannot be found",
		})
		return
	} else {
		response := Response{}
		copier.Copy(&response, &user)
		c.JSON(http.StatusOK, response)
	}
}

func (cr *UserHandler) LoginHandler(c *gin.Context) {
	var requestBody structs.LoginBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Internal Error"})
		return
	}
	if requestBody.Email == "" || requestBody.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Email or Password cannot be empty ",
			"Data":  requestBody,
		})
		return
	}

	user, err := cr.userUseCase.FindByEmail(c, requestBody.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Email or Password is incorrect",
			"Data":  requestBody,
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Email or Password is incorrect",
			"Data":  requestBody,
		})
		return
	}

	expirationTime := time.Now().Add(2 * time.Hour).Unix()
	claims := &structs.Claims{
		ID: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "It's not you, its us. Please try again later",
		})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"Error": false,
		"Data":  user,
	})
}

// should be implemented after the use of token
func (cr *UserHandler) UpdatePassword(c *gin.Context) {
	var requestBody structs.UpdatePassword

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Internal Error"})
		return
	}

	user, err := cr.userUseCase.FindByEmail(c, requestBody.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Email cannot be found",
			"Email": requestBody.Email,
		})
		return
	}

	hashedPassword, errmsg := requestBody.ValidateAndHash()
	if errmsg != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "password does not match",
		})
		return
	}
	user.Password = hashedPassword
	if _, err := cr.userUseCase.UpdatePassword(c, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "password does not match",
		})
		return
	}

	user.Password = ""
	c.JSON(http.StatusAccepted, gin.H{
		"Error": false,
		"Data":  user,
	})
}
