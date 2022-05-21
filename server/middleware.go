package main

import (
	"chat-app/server/db"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
