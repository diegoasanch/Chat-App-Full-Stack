package users

import (
	"chat-app/api/db"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type MyCustomClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}


func CreateUserToken(user *db.User) (string, error) {
	// Create claims while leaving out some of the optional fields
	createdAt := time.Now()
	expiresAt := createdAt.Add(time.Hour * 24)

	claims := MyCustomClaims{
		UserId:  user.BaseDbModel.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: createdAt.Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))

	return token.SignedString(secret)
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
