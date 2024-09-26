package database

import (
    "fmt"
    "log"
    "os"
    "github.com/Prototype-1/golang-auth-jwt/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
        log.Fatal("One or more required environment variables are not set")
    }

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    // Log successful connection
    log.Println("Successfully connected to the database")

    // Run migrations
    err = DB.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatalf("Error running migrations: %v", err)
    }

    // Log successful migration
    log.Println("Database migration completed")
}






