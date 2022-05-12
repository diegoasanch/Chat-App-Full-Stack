package db

import (
	"time"
)


type Message struct {
	ID string `gorm:"primary_key;default:gen_random_uuid();unique" json:"id"`

	Message string `json:"message"`
	UserId string `gorm:"not null" json:"user_id"`
	User User `gorm:"foreignKey:UserId;constraint:CnDelete:CASCADE" json:"user"`

	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
	DeletedAt time.Time `gorm:"index;default:null" json:"deleted_at"`
}

type User struct {
	ID string `gorm:"primary_key;default:gen_random_uuid();unique" json:"id"`

	Name string `json:"name"`
	Email string `gorm:"unique" json:"email"`
	Password string `json:"password"`

	Messages []Message

	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:null" json:"updated_at"`
	DeletedAt time.Time `gorm:"index;default:null" json:"deleted_at"`
}
