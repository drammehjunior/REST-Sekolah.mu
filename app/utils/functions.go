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

func HashPassword(password string) string {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(bs[:])
}

func IsLoginInputValid(body rest_structs.LoginBody) error {
	if body.Email == "" || body.Password == "" {
		return errors.New("email or password cannot be empty")
	}
	if _, err := mail.ParseAddress(body.Email); err != nil {
		return errors.New("email is not valid")
	}
	return nil
}

func IsEmailValid(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("email is not valid")
	}
	return nil
}

func IsPasswordMatched(oldPassword string, newPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(newPassword)); err != nil {
		return true
	}
	return false
}

func SignInToken(user domain.Users) (string, error) {

	expirationTime := time.Now().Add(2 * time.Hour).Unix()
	claims := &rest_structs.Claims{
		ID:    int64(user.Id),
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
