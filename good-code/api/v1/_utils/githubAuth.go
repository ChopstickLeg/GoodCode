package utils

import (
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
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

	privateKeyPEM := os.Getenv("GITHUB_APP_PRIVATE_KEY")
	if privateKeyPEM == "" {
		log.Println("ERROR: GITHUB_APP_PRIVATE_KEY environment variable not set")
		return "", 0, errors.New("GITHUB_APP_PRIVATE_KEY environment variable not set")
	}
	log.Printf("GitHub App Private Key found (length: %d)", len(privateKeyPEM))

	// Parse the private key
	log.Println("Parsing private key for RS256 signing")
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		log.Println("ERROR: Failed to parse PEM block containing the private key")
		return "", 0, errors.New("failed to parse PEM block containing the private key")
	}
	log.Printf("PEM block type: %s", block.Type)

	var privateKey *rsa.PrivateKey
	var err error

	if block.Type == "RSA PRIVATE KEY" {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else if block.Type == "PRIVATE KEY" {
		key, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
		if parseErr != nil {
			err = parseErr
		} else {
			var ok bool
			privateKey, ok = key.(*rsa.PrivateKey)
			if !ok {
				err = errors.New("parsed key is not an RSA private key")
			}
		}
	} else {
		err = fmt.Errorf("unsupported private key type: %s", block.Type)
	}

	if err != nil {
		log.Printf("ERROR: Failed to parse private key: %v", err)
		return "", 0, fmt.Errorf("failed to parse private key: %w", err)
	}
	log.Println("Successfully parsed private key")

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

	log.Println("Attempting to sign JWT with private key")
	jwtToken, err := token.SignedString(privateKey)
	if err != nil {
		log.Printf("ERROR: Failed to sign JWT token: %v", err)
		log.Printf("Error type: %T", err)
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
