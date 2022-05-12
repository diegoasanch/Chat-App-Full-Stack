package users

import (
	"chat-app/api/db"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func AuthMiddleware(c *gin.Context) {
	headerToken := c.GetHeader("Authorization")
	headerLen := len(headerToken)
	if headerLen < 7 || !strings.Contains(headerToken, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{ "status": "error", "message": "Invalid token" })
		c.Abort()
		return
	}
	tokenString := strings.Split(headerToken, "Bearer ")[1]

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{ "status": "error", "message": "No token provided" })
		c.Abort()
		return
	}

	token, err := VerifyUserToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{ "status": "error", "message": "Invalid token" })
		c.Abort()
		return
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{ "status": "error", "message": "Invalid token claims" })
		c.Abort()
		return
	}

	user := db.User{}
	result := db.DB.Where("id = ?", claims.UserId).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{ "status": "error", "message": "User does not exist" })
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}
