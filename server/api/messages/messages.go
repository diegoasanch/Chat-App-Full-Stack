package messages

import (
	"chat-app/server/api/users"
	"chat-app/server/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MessageRoutes(router *gin.RouterGroup) {
	router.Use(users.AuthMiddleware)
	router.POST("/", createMessage)
	router.GET("/", getMessages)
	router.POST("/delete", deleteMessage)
}

type CreateMessageBody struct {
	Message string `json:"message"`
}
func createMessage(c *gin.Context) {
	body := CreateMessageBody{}
	user := c.MustGet("user").(db.User)

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "status": "error", "message": "Invalid request" })
		return
	}
	message := db.Message{ Message: body.Message, UserId: user.ID }
	result := db.DB.Create(&message)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "status": "error", "message": message.ID })
		return
	}

	c.JSON(http.StatusOK, gin.H{ "status": "ok", "message": "Message created", "messageId": result.RowsAffected })
}

func getMessages(c *gin.Context) {
	var messages []db.Message
	db.DB.Preload("User").Model(&db.Message{}).Find(&messages).Limit(5).Order("created_at desc")

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

type DeleteMessageBody struct {
	ID uint `json:"id"`
}
func deleteMessage(c *gin.Context) {
	body := DeleteMessageBody{}


	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "message": "Invalid body" })
		return
	}
	result := db.DB.Where("id = ?", body.ID).Delete(&db.Message{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{ "error": "Invalid" })
		return
	}
	c.JSON(http.StatusOK, gin.H{ "message": "Deleted message", "rowsAffected": result.RowsAffected })
}
