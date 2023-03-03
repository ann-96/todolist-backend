package models

import "github.com/golang-jwt/jwt"

type Claims struct {
	UserID int `json:"UserID"`
	jwt.StandardClaims
}
