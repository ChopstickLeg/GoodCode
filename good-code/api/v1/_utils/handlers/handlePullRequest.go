package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"slices"

	db "github.com/chopstickleg/good-code/api/v1/_db"
	utils "github.com/chopstickleg/good-code/api/v1/_utils"

	"github.com/google/go-github/v72/github"
	"google.golang.org/genai"
)

var actions = []string{"opened", "synchronize", "reopened"}

func HandlePullRequestEvent(w http.ResponseWriter, body github.PullRequestEvent) {
	if slices.Contains(actions, body.GetAction()) {
		log.Printf("Received PR event: %s for PR #%d in %s", body.GetAction(), body.GetNumber(), body.GetRepo().GetFullName())
	} else {
		log.Printf("Received non-PR event: %s", body.GetAction())
		return
	}

	installationID := body.GetInstallation().GetID()
	if installationID == 0 {
		log.Printf("ERROR: No installation ID found in PR event")
		http.Error(w, "No installation ID found", http.StatusBadRequest)
		return
	}
	log.Printf("Using installation ID: %d", installationID)

	installationToken, err := utils.GetGitHubInstallationToken(installationID)
	if err != nil {
		log.Printf("Failed to get GitHub installation token: %v", err)
		http.Error(w, "Unable to get GitHub installation token", http.StatusInternalServerError)
		return
	}

	GHclient := github.NewClient(nil)
	authedGHClient := GHclient.WithAuthToken(installationToken)
	diff, _, err := authedGHClient.PullRequests.GetRaw(context.Background(), body.GetRepo().GetOwner().GetLogin(), body.GetRepo().GetName(), body.GetNumber(), github.RawOptions{
		Type: github.RawType(github.Diff),
	})
	if err != nil {
		log.Printf("Failed to get PR diff for PR #%d in %s: %v", body.GetNumber(), body.GetRepo().GetFullName(), err)
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
		log.Printf("Failed to create AI client: %v", err)
		http.Error(w, "Unable to create AI client", http.StatusInternalServerError)
		return
	}

	config := genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText("You are a code review assistant. You will be given a diff of a pull request. Your task is to review the code and provide feedback. You should be sarcastic and condescending, but still helpful and provide useful feedback that is factually accurate to the best of your knowledge", genai.RoleModel),
	}
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite-preview-06-17",
		genai.Text(diff),
		&config,
	)
	if err != nil {
		log.Printf("Failed to generate AI analysis: %v", err)
		http.Error(w, "Unable to generate AI analysis", http.StatusInternalServerError)
		return
	}

	if result == nil || result.Text() == "" {
		log.Printf("AI analysis returned empty result for PR #%d", body.GetNumber())
		http.Error(w, "AI analysis returned empty result", http.StatusInternalServerError)
		return
	}

	log.Println(result.Text())
	conn, err := db.GetDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
		return
	}

	pr := db.AiRoast{
		Content:           result.Text(),
		RepoID:            body.GetRepo().GetID(),
		PullRequestNumber: body.PullRequest.GetNumber(),
	}
	err = conn.Create(&pr).Error
	if err != nil {
		log.Printf("Failed to save AI analysis to database for PR #%d: %v", body.GetNumber(), err)
		http.Error(w, "Unable to save AI analysis to database", http.StatusInternalServerError)
		return
	}

	_, _, err = authedGHClient.PullRequests.CreateComment(context.Background(), body.GetRepo().GetOwner().GetLogin(), body.GetRepo().GetName(), body.GetNumber(), &github.PullRequestComment{
		Body: github.Ptr(result.Text()),
	})
	if err != nil {
		log.Printf("Failed to create comment on PR #%d in %s: %v", body.GetNumber(), body.GetRepo().GetFullName(), err)
		http.Error(w, "Unable to create comment on PR", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully processed PR #%d in %s", body.GetNumber(), body.GetRepo().GetFullName())

}
