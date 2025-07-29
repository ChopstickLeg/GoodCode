package utils

import (
	"context"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/v72/github"
)

func GenerateGitHubJWT() (string, error) {
	log.Println("Starting GitHub JWT generation")

	clientId := os.Getenv("GITHUB_APP_CLIENT_ID")
	if clientId == "" {
		log.Println("ERROR: GITHUB_APP_CLIENT_ID environment variable not set")
		return "", errors.New("GITHUB_APP_CLIENT_ID environment variable not set")
	}
	log.Printf("GitHub App Client ID found: %s", clientId)

	privateKeyPEM := os.Getenv("GITHUB_APP_PRIVATE_KEY")
	if privateKeyPEM == "" {
		log.Println("ERROR: GITHUB_APP_PRIVATE_KEY environment variable not set")
		return "", errors.New("GITHUB_APP_PRIVATE_KEY environment variable not set")
	}
	log.Printf("GitHub App Private Key found (length: %d)", len(privateKeyPEM))

	log.Println("Parsing private key for RS256 signing")
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		log.Println("ERROR: Failed to parse PEM block containing the private key")
		return "", errors.New("failed to parse PEM block containing the private key")
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
		return "", fmt.Errorf("failed to parse private key: %w", err)
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
		return "", fmt.Errorf("JWT signing failed: %w", err)
	}

	log.Printf("Successfully generated JWT token (length: %d)", len(jwtToken))
	return jwtToken, nil
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

func GetGitHubInstallationToken(installationID int64) (string, error) {
	log.Printf("Getting installation access token for installation ID: %d", installationID)

	jwt, err := GenerateGitHubJWT()
	if err != nil {
		log.Printf("ERROR: Failed to get GitHub JWT: %v", err)
		return "", fmt.Errorf("failed to get GitHub JWT: %w", err)
	}

	client := github.NewClient(nil).WithAuthToken(jwt)

	installationToken, _, err := client.Apps.CreateInstallationToken(
		context.Background(),
		installationID,
		&github.InstallationTokenOptions{},
	)
	if err != nil {
		log.Printf("ERROR: Failed to create installation token for installation %d: %v", installationID, err)
		return "", fmt.Errorf("failed to create installation token: %w", err)
	}

	log.Printf("Successfully obtained installation access token (expires at: %s)", installationToken.GetExpiresAt().Format(time.RFC3339))
	return installationToken.GetToken(), nil
}
