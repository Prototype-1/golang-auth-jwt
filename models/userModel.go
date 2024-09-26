package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID                      uint           `gorm:"primaryKey" json:"id"`
    FirstName         string         `json:"first_name" validate:"required,min=2,max=100"`
    LastName         string         `json:"last_name" validate:"required,min=2,max=100"`
    Password          string         `json:"password" validate:"required,min=6"`
    Email                 string         `json:"email" validate:"required,email"`
    Phone               string         `json:"phone" validate:"required"`
    Token                string         `json:"token"`
    UserType          string         `json:"user_type"`
    RefreshToken   string         `json:"refresh_token"`
    CreatedAt        time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    UserID             string         `json:"user_id"`

     // Adding GORM-specific field for soft delete
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

   




