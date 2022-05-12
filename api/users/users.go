package users

import (
	"chat-app/api/db"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func UserRoutes(router *gin.RouterGroup) {
	router.POST("/signup", createUser)
	router.POST("/login", login)
}

type CreateUserBody struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}
func createUser(c *gin.Context) {
	body := CreateUserBody{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "status": "error" })
		return
	}
	if body.Name == "" || body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{ "status": "error", "message": "Invalid payload" })
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "status": "error", "message": "Error hashing password" })
		return
	}

	user := &db.User{ Name: body.Name, Email: body.Email, Password: string(hash) }

	result := db.DB.Create(user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "status": "error", "message": result.Error.Error() })
		return
	}
	c.JSON(http.StatusOK, gin.H{ "status": "ok", "message": "User created", "userId":  user.BaseDbModel.ID  })
}

type LoginBody struct {
	Email string `json:"email"`
	Password string `json:"password"`
}
func login(c *gin.Context) {
	body := LoginBody{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "status": "error" })
		return
	}
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{ "status": "error", "message": "Invalid payload" })
		return
	}
	var user db.User
	result := db.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "status": "error", "message": result.Error.Error() })
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{ "status": "error", "message": "Invalid credentials" })
		return
	}

	hmacSampleSecret := []byte("secretCode")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.BaseDbModel.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "status": "error", "message": "Error signing token" })
		fmt.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{ "status": "ok", "message": "User logged in", "token": tokenString })
}
