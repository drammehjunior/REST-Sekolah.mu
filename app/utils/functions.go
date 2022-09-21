package utils

import (
	"errors"
	"exampleclean.com/refactor/app/domain"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"time"
)

var jwtKey = []byte("my_super_secret_key")

func IsLoginInputValid(body rest_structs.LoginBody) error {

	if body.Email == "" || body.Password == "" {
		return errors.New("email or password cannot be empty")
	}
	if _, err := mail.ParseAddress(body.Email); err != nil {
		return errors.New("email is not valid")
	}
	return nil
}
func IsPasswordMatched(oldPassword string, newPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(newPassword)); err != nil {
		return false
	}
	return true
}

func SignInToken(user domain.Users) (string, error) {

	expirationTime := time.Now().Add(2 * time.Hour).Unix()
	claims := &rest_structs.Claims{
		ID:    user.Id,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if tokenString, err := token.SignedString(jwtKey); err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}
