package users

import (
	"chat-app/server/db"
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

	claims := &MyCustomClaims{
		UserId:  user.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: createdAt.Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))

	return token.SignedString(secret)
}

