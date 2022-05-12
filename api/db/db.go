package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbPort := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s",dbHost, dbUser, dbPassword, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = db
	migrateDB()
}

func migrateDB() {
	dbModels := []interface{}{
		&Message{},
		&User{},
	}

	for _, model := range dbModels {
		if err := DB.AutoMigrate(model); err != nil {
			panic("failed to migrate database:\n" + err.Error())
		}
	}

}
