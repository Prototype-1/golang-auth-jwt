package main

import (
    //"html/template"
    //"net/http"
    "log"
    "os"
    "time"
    "github.com/Prototype-1/golang-auth-jwt/database"
    "github.com/Prototype-1/golang-auth-jwt/routes"
   // "github.com/Prototype-1/golang-auth-jwt/models"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
   // "github.com/dgrijalva/jwt-go"
)

func NoCache() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
        c.Header("Cache-Control", "post-check=0, pre-check=0")
        c.Header("Pragma", "no-cache")
        c.Next()
    }
}

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    database.Connect()

    router := gin.Default()

    // Middleware setup
    router.Use(gin.Logger())
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:8000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Static files
    router.Static("/frontend", "./frontend")
    router.LoadHTMLGlob("frontend/*")


    // HTML Routes
    router.GET("/", func(c *gin.Context) { c.HTML(200, "index.html", nil) })
    router.GET("/admin/signup", func(c *gin.Context) { c.HTML(200, "admin_signup.html", nil) })
    router.GET("/admin/login", func(c *gin.Context) { c.HTML(200, "admin_login.html", nil) })
    router.GET("/users/signup", func(c *gin.Context) { c.HTML(200, "user_signup.html", nil) })
    router.GET("/users/login", func(c *gin.Context) { c.HTML(200, "user_login.html", nil) })


    // Register route groups
    routes.AdminRoutes(router)
    routes.AuthRoutes(router)
    routes.UserRoutes(router)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }
    log.Printf("Starting server on port %s...", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}













