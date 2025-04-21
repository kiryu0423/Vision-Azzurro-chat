package util

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() []byte {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		panic("JWT_SECRET is not set")
	}
	return []byte(key)
}

func GenerateJWT(userID uint, userName string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": userName,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

func ValidateJWT(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return getSecretKey(), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}
	claims := token.Claims.(jwt.MapClaims)
	return uint(claims["user_id"].(float64)), nil
}

func ValidateJWTAndExtract(tokenStr string) (uint, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return getSecretKey(), nil
	})
	if err != nil || !token.Valid {
		return 0, "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	userIDFloat, ok1 := claims["user_id"].(float64)
	userName, ok2 := claims["user_name"].(string)
	if !ok1 || !ok2 {
		return 0, "", err
	}
	return uint(userIDFloat), userName, nil
}
