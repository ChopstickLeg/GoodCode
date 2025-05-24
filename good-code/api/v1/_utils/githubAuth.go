package utils

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	cachedJWT      string
	cachedJWTExp   int64
	cachedJWTMutex sync.Mutex
)

func GenerateGitHubJWT() (string, int64, error) {
	clientId := os.Getenv("GITHUB_APP_CLIENT_ID")
	if clientId == "" {
		return "", 0, errors.New("GITHUB_APP_CLIENT_ID environment variable not set")
	}
	clientSecret := os.Getenv("GITHUB_APP_CLIENT_SECRET")
	if clientSecret == "" {
		return "", 0, errors.New("GITHUB_APP_CLIENT_SECRET environment variable not set")
	}
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().Add(-time.Minute).Unix()
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	claims["iss"] = clientId
	claims["alg"] = "RS256"

	jwtToken, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return "", 0, err
	}
	return jwtToken, claims["exp"].(int64), nil
}
func GetGitHubJWT() (string, error) {
	cachedJWTMutex.Lock()
	defer cachedJWTMutex.Unlock()

	now := time.Now().Unix()
	if cachedJWT != "" && cachedJWTExp > now {
		return cachedJWT, nil
	}
	token, exp, err := GenerateGitHubJWT()
	if err != nil {
		return "", err
	}
	cachedJWT = token
	cachedJWTExp = exp
	return cachedJWT, nil
}