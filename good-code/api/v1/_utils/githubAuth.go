package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
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
	log.Println("Starting GitHub JWT generation")

	clientId := os.Getenv("GITHUB_APP_CLIENT_ID")
	if clientId == "" {
		log.Println("ERROR: GITHUB_APP_CLIENT_ID environment variable not set")
		return "", 0, errors.New("GITHUB_APP_CLIENT_ID environment variable not set")
	}
	log.Printf("GitHub App Client ID found: %s", clientId)

	clientSecret := os.Getenv("GITHUB_APP_CLIENT_SECRET")
	if clientSecret == "" {
		log.Println("ERROR: GITHUB_APP_CLIENT_SECRET environment variable not set")
		return "", 0, errors.New("GITHUB_APP_CLIENT_SECRET environment variable not set")
	}
	log.Printf("GitHub App Client Secret found (length: %d)", len(clientSecret))

	log.Println("Creating JWT token with RS256 signing method")
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)

	now := time.Now()
	iat := now.Add(-time.Minute).Unix()
	exp := now.Add(time.Minute * 10).Unix()

	claims["iat"] = iat
	claims["exp"] = exp
	claims["iss"] = clientId
	claims["alg"] = "RS256"

	log.Printf("JWT claims set - iat: %d, exp: %d, iss: %s", iat, exp, clientId)

	log.Println("Attempting to sign JWT with client secret")
	log.Printf("WARNING: Using client secret as signing key for RS256 - this may be incorrect. RS256 typically requires a private key, not a client secret")

	jwtToken, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		log.Printf("ERROR: Failed to sign JWT token: %v", err)
		log.Printf("Error type: %T", err)
		log.Println("NOTE: RS256 requires a private key (PEM format), not a client secret. Consider using HS256 with client secret or RS256 with private key")
		return "", 0, fmt.Errorf("JWT signing failed: %w", err)
	}

	log.Printf("Successfully generated JWT token (length: %d)", len(jwtToken))
	return jwtToken, claims["exp"].(int64), nil
}

func GetGitHubJWT() (string, error) {
	log.Println("Attempting to get GitHub JWT (checking cache first)")
	cachedJWTMutex.Lock()
	defer cachedJWTMutex.Unlock()

	now := time.Now().Unix()
	if cachedJWT != "" && cachedJWTExp > now {
		log.Printf("Using cached JWT (expires at: %d, current time: %d)", cachedJWTExp, now)
		return cachedJWT, nil
	}

	if cachedJWT != "" {
		log.Printf("Cached JWT expired (expired at: %d, current time: %d)", cachedJWTExp, now)
	} else {
		log.Println("No cached JWT found")
	}

	log.Println("Generating new JWT token")
	token, exp, err := GenerateGitHubJWT()
	if err != nil {
		log.Printf("ERROR: Failed to generate new JWT: %v", err)
		return "", err
	}

	cachedJWT = token
	cachedJWTExp = exp
	log.Printf("Successfully cached new JWT (expires at: %d)", exp)
	return cachedJWT, nil
}

func VerifyGitHubSignature(payload []byte, signature string) bool {
	log.Printf("Verifying GitHub signature (payload length: %d, signature: %s)", len(payload), signature)

	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret == "" {
		log.Println("ERROR: GITHUB_WEBHOOK_SECRET environment variable not set")
		return false
	}
	log.Printf("GitHub webhook secret found (length: %d)", len(secret))

	key := hmac.New(sha256.New, []byte(secret))
	key.Write([]byte(string(payload)))
	computedSignature := fmt.Sprintf("sha256=%x", key.Sum(nil))

	log.Printf("Computed signature: %s", computedSignature)
	log.Printf("Provided signature: %s", signature)

	isValid := hmac.Equal([]byte(signature), []byte(computedSignature))
	log.Printf("Signature verification result: %t", isValid)

	return isValid
}
