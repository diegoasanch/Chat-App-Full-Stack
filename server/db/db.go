package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global DB connection
var DB *gorm.DB

const MAX_CONNECTION_TRIES = 10
const INIT_ATTEMPT_COUNT = 0

func Initialize() {
	ConnectDB(INIT_ATTEMPT_COUNT)
}

// Initial attempt should be 0
func ConnectDB(attempt int) {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbPort := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s",dbHost, dbUser, dbPassword, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{  })

	// Retry up to MAX_CONNECTION_TRIES times
	if err != nil {
		fmt.Printf("attempt: %d, failed to connect to database:\n", attempt)
		if attempt < MAX_CONNECTION_TRIES {
			time.Sleep(500 * time.Millisecond)
			ConnectDB(attempt + 1)
		} else {
			fmt.Println("failed to connect to database", err.Error())
			panic("failed to connect database")
		}
	}

	migrateDB(db)
	DB = db
}

func migrateDB(db *gorm.DB) {
	dbModels := []interface{}{
		&Message{},
		&User{},
	}

	for _, model := range dbModels {
		if err := db.AutoMigrate(model); err != nil {
			panic("failed to migrate database:\n" + err.Error())
		}
	}

}
