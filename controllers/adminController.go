package controllers

import (
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "net/http"
    "os"
    "log"
    "fmt"
    "time"
    "github.com/Prototype-1/golang-auth-jwt/models"
    "github.com/Prototype-1/golang-auth-jwt/utils"
    "github.com/Prototype-1/golang-auth-jwt/database"
    "golang.org/x/crypto/bcrypt"
     "github.com/google/uuid"
)

// LoginInput represents the expected input for login
type LoginInput struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// SignupInput represents the expected input for signup
type SignupInput struct {
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    Email     string `json:"email" binding:"required"`
    Password  string `json:"password" binding:"required"`
    Phone     string `json:"phone" binding:"required"`
    UserType  string `json:"user_type"`
}

// AdminSignupInput represents the expected input for admin signup
type AdminSignupInput struct {
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    Email     string `json:"email" binding:"required"`
    Password  string `json:"password" binding:"required"`
    Phone     string `json:"phone" binding:"required"`
}

func generateTokens(user *models.User) (string, error) {
    claims := jwt.MapClaims{}
    claims["authorized"] = true
    claims["user_id"] = user.UserID
    claims["role"] = user.UserType // Add role to claims
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func AdminSignup(c *gin.Context) {
    var input AdminSignupInput
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Printf("Admin creation failed: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if email already exists
    var existingAdmin models.User
    if err := database.DB.Where("email = ?", input.Email).First(&existingAdmin).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }

    admin := models.User {
        FirstName: input.FirstName,
        LastName:  input.LastName,
        Email:     input.Email,
        Password:  string(hashedPassword),
        Phone:     input.Phone,
        UserType:  "admin",
        UserID:    uuid.New().String(),
    }

    if err := database.DB.Create(&admin).Error; err != nil {
        log.Printf("Failed to create admin: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Admin created successfully"})
}


func AdminLogIn(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var admin models.User
    if err := database.DB.Where("email = ?", input.Email).First(&admin).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    if !utils.CheckPasswordHash(input.Password, admin.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Check if the user is an admin
    if admin.UserType != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
        return
    }

    token, err := generateToken(&admin)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}


// GetAllUsers retrieves all users
func GetAllUsers(c *gin.Context) {
    role, exists := c.Get("role")
    if !exists || role != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
        return
    }

    var users []models.User
    result := database.DB.Find(&users)
    if result.Error != nil {
        log.Printf("Failed to retrieve users: %v", result.Error)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
        return
    }
    c.JSON(http.StatusOK, users)
}

// Create a new user
func CreateUser(c *gin.Context) {
    var input SignupInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fmt.Println("Signup input received:", input)

    var existingUser models.User
    if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
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


// GetUser retrieves a single user by ID
func GetUser(c *gin.Context) {
    var user models.User
    if err := database.DB.Where("user_id = ?", c.Param("user_id")).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

// UpdateUser updates a user's information
func UpdateUser(c *gin.Context) {
    userID := c.Param("id")

    var input models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Check if the email is already taken by another user
    if input.Email != "" && input.Email != user.Email {
        var existingUser models.User
        if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use by another user"})
            return
        }
    }

    // Update the user's fields only if they have been changed
    if input.FirstName != "" {
        user.FirstName = input.FirstName
    }
    if input.LastName != "" {
        user.LastName = input.LastName
    }
    if input.Email != "" {
        user.Email = input.Email
    }
    if input.Phone != "" {
        user.Phone = input.Phone
    }
    if input.UserType != "" {
        user.UserType = input.UserType
    }
    if input.Password != "" {
        hashedPassword, _ := utils.HashPassword(input.Password) // Assume HashPassword is a function that hashes passwords
        user.Password = string(hashedPassword)
    }

    if err := database.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}



// DeleteUser deletes a user by ID
func DeleteUser(c *gin.Context) {
    userID := c.Param("id")

    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if err := database.DB.Delete(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}


