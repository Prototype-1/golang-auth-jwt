package utils

import (
    "github.com/dgrijalva/jwt-go"
    "fmt"
    "os"
    "time"
)

type Claims struct {
    UserID   uint   `json:"user_id"`
    Role     string `json:"role"`
    jwt.StandardClaims
}

func GenerateToken(userID uint, role string) (string, error) {
    claims := &Claims{
        UserID: userID,
        Role:   role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        if ve, ok := err.(*jwt.ValidationError); ok {
            if ve.Errors&jwt.ValidationErrorMalformed != 0 {
                return nil, fmt.Errorf("token is malformed")
            } else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
                return nil, fmt.Errorf("token is expired or not yet valid")
            } else {
                return nil, fmt.Errorf("couldn't handle this token: %v", err)
            }
        }
    }

    return token, nil
}