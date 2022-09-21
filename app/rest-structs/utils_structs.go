package rest_structs

import "github.com/golang-jwt/jwt"

type Claims struct {
	ID    int64  `json:"ID"`
	Email string `json:"Email"`
	jwt.StandardClaims
}
