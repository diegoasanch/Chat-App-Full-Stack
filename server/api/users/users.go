package users

import (
	"chat-app/server/db"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, gin.H{ "status": "ok", "message": "User created", "userId":  user.ID  })
}

type LoginBody struct {
	Email string `json:"email"`
	Password string `json:"password"`
}
func login(c *gin.Context) {
	body := LoginBody{}
	if err := c.BindJSON(&body); err != nil || body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{ "status": "error", "message": "Invalid payload" })
		return
	}

	var user db.User
	result := db.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "status": "error", "message": result.Error.Error() })
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{ "status": "error", "message": "Incorrect password" })
		return
	}

	token, err := CreateUserToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "status": "error", "message": "Error Creating token" })
		fmt.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{ "status": "ok", "message": "User logged in", "token": token })
}
