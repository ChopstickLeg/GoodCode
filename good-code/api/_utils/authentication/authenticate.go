package authentication

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func VerifyToken(tokenString string) error {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return fmt.Errorf("JWT_SECRET_KEY environment variable not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
