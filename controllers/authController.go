package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/Prototype-1/golang-auth-jwt/models"
    "github.com/Prototype-1/golang-auth-jwt/database"
    "golang.org/x/crypto/bcrypt"
    "net/http"
    "time"
    "os"
    //"log"
    "fmt"
    "github.com/google/uuid"
    "github.com/dgrijalva/jwt-go"
)

type LogIninput struct {
    Email        string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type SignUpinput struct {
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    Email     string `json:"email" binding:"required"`
    Password  string `json:"password" binding:"required"`
    Phone     string `json:"phone" binding:"required"`
    UserType  string `json:"user_type"`
}

func generateToken(user *models.User) (string, error) {
    claims := jwt.MapClaims{}
    claims["authorized"] = true
    claims["user_id"] = user.UserID
    claims["role"] = user.UserType
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}


func SignUp(c *gin.Context) {
    var input SignUpinput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fmt.Println("Signup input received:", input)

    // Check if the email already exists
    var existingUser models.User
    if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
        return
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user := models.User{
        FirstName: input.FirstName,
        LastName:  input.LastName,
        Email:     input.Email,
        Password:  string(hashedPassword),
        Phone:     input.Phone,
        UserType:  input.UserType,
        UserID:    uuid.New().String(), 
    }

    fmt.Println("Attempting to create user:", user)

    if err := database.DB.Create(&user).Error; err != nil {
        fmt.Println("Error creating user:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    fmt.Println("User created:", user)
    c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}


func Login(c *gin.Context) {
    var input LogIninput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    result := database.DB.Where("email = ?", input.Email).First(&user)
    if result.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := generateTokens(&user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    user.Token = token
    database.DB.Save(&user)

    c.JSON(http.StatusOK, gin.H{
        "token": token,
    "username": user.Email,
})
}









   





