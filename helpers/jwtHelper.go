package helpers

import (
    "time"
    "os"
    "github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
    Email     string
    FirstName string
    LastName  string
    UserType  string
    jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("JWT_SECRET")

func GenerateAllTokens(email, firstname, lastname, userType string) (signedToken string, signedRefreshToken string, err error) {
    claims := &SignedDetails{
        Email:     email,
        FirstName: firstname,
        LastName:  lastname,
        UserType:  userType,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
        },
    }

    refreshClaims := &SignedDetails{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
        },
    }

    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
    if err != nil {
        return "", "", err
    }

    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
    if err != nil {
        return "", "", err
    }

    return token, refreshToken, nil
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
    token, err := jwt.ParseWithClaims(
        signedToken,
        &SignedDetails{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(SECRET_KEY), nil
        },
    )

    if err != nil {
        msg = err.Error()
        return
    }

    claims, ok := token.Claims.(*SignedDetails)
    if !ok {
        msg = "the token is invalid"
        return
    }

    if claims.ExpiresAt < time.Now().Local().Unix() {
        msg = "token is expired"
        return
    }

    return claims, msg
}
