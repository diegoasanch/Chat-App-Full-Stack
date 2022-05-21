package main

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

type MyCustomClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}


func VerifyUserToken(tokenString string) (*jwt.Token, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, err
}

