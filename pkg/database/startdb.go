package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func StartDB() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("cannot load env")
	}
	dsn := os.Getenv("DB_DSN")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("faild to connect to database")
	}
}
