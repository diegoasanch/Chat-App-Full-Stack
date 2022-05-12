package db

import (
	"time"
)

type BaseDbModel struct {
	ID string `gorm:"primary_key;default:gen_random_uuid();unique" json:"id"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
	DeletedAt time.Time `gorm:"index;default:null" json:"deleted_at"`
}

type Message struct {
	BaseDbModel BaseDbModel `gorm:"embedded"`
	Message string `json:"message"`
	UserId string `json:"user_id"`
	User User `gorm:"foreignKey:UserId;constraint:CnDelete:CASCADE" json:"user"`
}

type User struct {
	BaseDbModel BaseDbModel `gorm:"embedded"`
	Name string `json:"name"`
	Email string `gorm:"unique" json:"email"`
	Password string `json:"password"`

	Messages []Message
}
