package middleware

import (
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func GetUserIDFromJWT(tokenString string) (int, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Printf("Error: JWT_SECRET_KEY environment variable not set")
		return 0, fmt.Errorf("JWT_SECRET_KEY environment variable not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {

		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("could not parse claims")
	}

	userID, ok := claims["id"].(int64)
	if !ok {
		return 0, fmt.Errorf("user ID not found in token")
	}

	return userID, nil
}
