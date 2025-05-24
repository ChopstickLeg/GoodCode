package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	middleware "github.com/chopstickleg/good-code/api/v1/_middleware"
	utils "github.com/chopstickleg/good-code/api/v1/_utils"

	"github.com/google/go-github/v72/github"
	"google.golang.org/genai"
)

var actions = []string{"opened", "synchronize", "reopened"}

func Handler(w http.ResponseWriter, r *http.Request) {
	middleware.AllowMethods(http.MethodPost)(AddPRHandler)(w, r)
}

func AddPRHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	if !VerifyGitHubSignature(body, r.Header.Get("X-Hub-Signature-256")) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	var requestBody github.PullRequestEvent
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Unable to decode GitHub event", http.StatusBadRequest)
		return
	}

	if slices.Contains(actions, requestBody.GetAction()) {
		fmt.Println("Received PR event:", requestBody.GetAction())
	} else {
		fmt.Println("Received non-PR event:", requestBody.GetAction())
		return
	}

	githubJWT, err := utils.GetGitHubJWT()
	if err != nil {
		http.Error(w, "Unable to get GitHub JWT", http.StatusInternalServerError)
		return
	}

	GHclient := github.NewClient(nil)
	authedGHClient := GHclient.WithAuthToken(githubJWT)
	diff, _, err := authedGHClient.PullRequests.GetRaw(context.Background(), requestBody.GetRepo().GetOwner().GetLogin(), requestBody.GetRepo().GetName(), requestBody.GetNumber(), github.RawOptions{
		Type: github.RawType(github.Diff),
	})

	if err != nil {
		http.Error(w, "Unable to get PR diff", http.StatusInternalServerError)
		return
	}

	token := os.Getenv("AI_API_TOKEN")
	if token == "" {
		http.Error(w, "Unable to get AI API token", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  token,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		http.Error(w, "Unable to create AI client", http.StatusInternalServerError)
		return
	}

	config := genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText("You are a code review assistant. You will be given a diff of a pull request. Your task is to review the code and provide feedback. You should be sarcastic and condescending, but still helpful and provide useful feedback that is factually accurate to the best of your knowledge", genai.RoleModel),
	}

	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-lite",
		genai.Text(diff),
		&config,
	)

	fmt.Print(result.Text())
}

func VerifyGitHubSignature(payload []byte, signature string) bool {
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	key := hmac.New(sha256.New, []byte(secret))
	key.Write([]byte(string(payload)))
	computedSignature := fmt.Sprintf("sha256=%x", key.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(computedSignature))
}
