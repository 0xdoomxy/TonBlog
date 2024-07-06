package utils

import "github.com/golang-jwt/jwt"

type JwtCustomClaims struct {
	jwt.StandardClaims
	PublicKey string `json:"publickey"`
	Address   string `json:"address"`
}
