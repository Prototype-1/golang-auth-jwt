package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/Prototype-1/golang-auth-jwt/database"
    "github.com/Prototype-1/golang-auth-jwt/models"
)

// GetUser function to get user by ID
func GetUsers(c *gin.Context) {
    userId := c.Param("user_id")

    var user models.User
    if err := database.DB.First(&user, "id = ?", userId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found, please contact admin"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": user})
}



