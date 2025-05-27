package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"slices"

	db "github.com/chopstickleg/good-code/api/v1/_db"

	"github.com/google/go-github/v72/github"
	"google.golang.org/genai"
)

var actions = []string{"opened", "synchronize", "reopened"}

func AddPRHandler(w http.ResponseWriter, body github.PullRequestEvent) {
	if slices.Contains(actions, body.GetAction()) {
		fmt.Println("Received PR event:", body.GetAction())
	} else {
		fmt.Println("Received non-PR event:", body.GetAction())
		return
	}

	githubJWT, err := GetGitHubJWT()
	if err != nil {
		http.Error(w, "Unable to get GitHub JWT", http.StatusInternalServerError)
		return
	}

	GHclient := github.NewClient(nil)
	authedGHClient := GHclient.WithAuthToken(githubJWT)
	diff, _, err := authedGHClient.PullRequests.GetRaw(context.Background(), body.GetRepo().GetOwner().GetLogin(), body.GetRepo().GetName(), body.GetNumber(), github.RawOptions{
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

	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
		return
	}

	pr := db.AiRoast{
		AiAnalysis:    result.Text(),
		RepoId:        body.GetRepo().GetID(),
		PullRequestId: body.PullRequest.GetID(),
	}

	err = conn.Create(&pr).Error
	if err != nil {
		http.Error(w, "Unable to save AI analysis to database", http.StatusInternalServerError)
		return
	}

	_, _, err = authedGHClient.PullRequests.CreateComment(context.Background(), body.GetRepo().GetOwner().GetLogin(), requestBody.GetRepo().GetName(), requestBody.GetNumber(), &github.PullRequestComment{
		Body: github.Ptr(result.Text()),
	})

	if err != nil {
		http.Error(w, "Unable to create comment on PR", http.StatusInternalServerError)
		return
	}

}

func VerifyGitHubSignature(payload []byte, signature string) bool {
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	key := hmac.New(sha256.New, []byte(secret))
	key.Write([]byte(string(payload)))
	computedSignature := fmt.Sprintf("sha256=%x", key.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(computedSignature))
}
