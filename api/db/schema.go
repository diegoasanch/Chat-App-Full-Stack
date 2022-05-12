package db

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Message string `json:"message"`
}
