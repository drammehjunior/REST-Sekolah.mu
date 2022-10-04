package utils

import (
	"crypto/sha1"
	"errors"
	"exampleclean.com/refactor/app/domain"
	rest_structs "exampleclean.com/refactor/app/rest-structs"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"time"
)

var jwtKey = []byte("my_super_secret_key")

func HashThisSHA1(password string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte("hash_salt"))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}

func HashPassword(password string) string {
	fmt.Println(viper.Get("ENCRYPT_KEY"))
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(bs[:])
}

func IsPasswordFilledAndMatched(first, second string) error {
	if first == "" || second == "" {
		return errors.New("passwords cannot be empty")
	}

	if first != second {
		return errors.New("password do not match")
	}

	return nil
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
	newPasswordHashed := HashThisSHA1(newPassword)
	if newPasswordHashed != oldPassword {
		return false
	} else {

		return true
	}
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
