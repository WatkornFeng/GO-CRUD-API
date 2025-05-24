// database/db.go
package database

import (
	"fmt"
	"log"
	"os"
	"project_restfulApi_go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host = "localhost" // or the Docker service name if running in another container
	port = 5432        // default PostgreSQL port
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("Database migration completed!")
	DB = db

}
