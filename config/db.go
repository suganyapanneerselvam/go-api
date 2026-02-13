package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		dsn = "host=localhost user=postgres password=Suganya@17 dbname=goapi port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}

	DB = db

	log.Println("✅ Database connected")
}
