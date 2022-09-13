package rest_structs

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type RequestSignup struct {
	Email           string `copier:"must"`
	Password        string `copier:"must"`
	PasswordConfirm string `copier:"must"`
	Firstname       string `copier:"must"`
	Lastname        string `copier:"must"`
}

type LoginBody struct {
	Email    string `copier:"must"`
	Password string `copier:"must"`
}

type UpdatePassword struct {
	Email              string `copier:"must"`
	OldPasswrod        string `copier:"must"`
	NewPassword        string `copier:"must"`
	NewPasswordConfirm string `copier:"must"`
}

func (c RequestSignup) ValidateAndHash() (string, error) {
	if c.Password != c.PasswordConfirm {
		return "", errors.New("password does not match")
	} else {
		bs, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.MinCost)
		if err != nil {
			fmt.Println(err)
		}
		return string(bs[:]), nil
	}
}

func (c UpdatePassword) ValidateAndHash() (string, error) {
	if c.NewPassword != c.NewPasswordConfirm {
		return "", errors.New("password does not match")
	} else {
		bs, err := bcrypt.GenerateFromPassword([]byte(c.NewPassword), bcrypt.MinCost)
		if err != nil {
			fmt.Println(err)
		}
		return string(bs[:]), nil
	}
}
