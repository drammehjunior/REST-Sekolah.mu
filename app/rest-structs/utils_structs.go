package rest_structs

import "github.com/golang-jwt/jwt"

type Claims struct {
	ID    int64  `json:"ID"`
	Email string `json:"Email"`
	jwt.StandardClaims
}

type RequestSignup struct {
	Email           string `copier:"must"`
	Password        string `copier:"must"`
	PasswordConfirm string `copier:"must"`
	Firstname       string `copier:"must"`
	Lastname        string `copier:"must"`
}

type SignupResponse struct {
	Email     string `copier:"must"`
	Firstname string `copier:"must"`
	Lastname  string `copier:"must"`
}
