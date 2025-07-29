package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"

	db "github.com/chopstickleg/good-code/api/_db"
	utils "github.com/chopstickleg/good-code/api/_utils"

	"github.com/google/go-github/v72/github"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

var actions = []string{"opened", "synchronize", "reopened"}

func HandlePullRequestEvent(w http.ResponseWriter, body github.PullRequestEvent) {
	conn, err := db.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	if slices.Contains(actions, body.GetAction()) {
		log.Printf("Received PR event: %s for PR #%d in %s", body.GetAction(), body.GetNumber(), body.GetRepo().GetFullName())
		if err := roastPullRequest(conn, &body); err != nil {
			log.Printf("Failed to roast PR #%d in %s: %v", body.GetNumber(), body.GetRepo().GetFullName(), err)
			http.Error(w, "Failed to roast pull request", http.StatusInternalServerError)
			return
		}
	} else if body.GetAction() == "closed" && body.PullRequest.GetMerged() {
		log.Printf("Received merged/closed PR event: %s for PR #%d in %s", body.GetAction(), body.GetNumber(), body.GetRepo().GetFullName())
	} else {
		log.Printf("Received non-PR event: %s", body.GetAction())
		return
	}
}

func getAuthedClient(installationID int64) (*github.Client, error) {
	if installationID == 0 {
		log.Printf("ERROR: No installation ID found in PR event")
		return nil, fmt.Errorf("No installation ID found")
	}
	log.Printf("Using installation ID: %d", installationID)

	installationToken, err := utils.GetGitHubInstallationToken(installationID)
	if err != nil {
		log.Printf("Failed to get GitHub installation token: %v", err)
		return nil, fmt.Errorf("Unable to get GitHub installation token: %w", err)
	}

	GHclient := github.NewClient(nil)
	authedGHClient := GHclient.WithAuthToken(installationToken)
	return authedGHClient, nil
}

func roastPullRequest(conn *gorm.DB, body *github.PullRequestEvent) error {
	authedGHClient, err := getAuthedClient(body.GetInstallation().GetID())
	if err != nil {
		log.Printf("Failed to get authenticated GitHub client: %v", err)
		return fmt.Errorf("Failed to get authenticated GitHub client: %w", err)
	}
	diff, _, err := authedGHClient.PullRequests.GetRaw(context.Background(), body.GetRepo().GetOwner().GetLogin(), body.GetRepo().GetName(), body.GetNumber(), github.RawOptions{
		Type: github.RawType(github.Diff),
	})
	if err != nil {
		log.Printf("Failed to get PR diff for PR #%d in %s: %v", body.GetNumber(), body.GetRepo().GetFullName(), err)
		return fmt.Errorf("Failed to get PR diff: %w", err)
	}

	token := os.Getenv("AI_API_TOKEN")
	if token == "" {
		log.Printf("ERROR: Unable to get AI API token")
		return fmt.Errorf("Unable to get AI API token")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  token,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Printf("Failed to create AI client: %v", err)
		return fmt.Errorf("Unable to create AI client: %w", err)
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
		return fmt.Errorf("Unable to generate AI analysis: %w", err)
	}

	if result == nil || result.Text() == "" {
		log.Printf("AI analysis returned empty result for PR #%d", body.GetNumber())
		return fmt.Errorf("AI analysis returned empty result")
	}

	log.Println(result.Text())

	pr := db.AiRoast{
		Content:           result.Text(),
		RepoID:            body.GetRepo().GetID(),
		PullRequestNumber: body.PullRequest.GetNumber(),
	}
	err = conn.Create(&pr).Error
	if err != nil {
		log.Printf("Failed to save AI analysis to database for PR #%d: %v", body.GetNumber(), err)
		return fmt.Errorf("Unable to save AI analysis to database: %w", err)
	}

	_, _, err = authedGHClient.Issues.CreateComment(context.Background(), body.GetRepo().GetOwner().GetLogin(), body.GetRepo().GetName(), body.GetNumber(), &github.IssueComment{
		Body: github.Ptr(result.Text()),
	})
	if err != nil {
		log.Printf("Failed to create comment on PR #%d in %s: %v", body.GetNumber(), body.GetRepo().GetFullName(), err)
		return fmt.Errorf("Unable to create comment on PR: %w", err)
	}

	log.Printf("Successfully processed PR #%d in %s", body.GetNumber(), body.GetRepo().GetFullName())
	return nil
}

func disablePullRequest(conn *gorm.DB, body *github.PullRequestEvent) error {

	err := conn.Model(&db.AiRoast{}).
		Where(&db.AiRoast{RepoID: body.GetRepo().GetID(), PullRequestNumber: body.GetNumber()}).
		Updates(&db.AiRoast{IsOpen: false}).
		Error
	if err != nil {
		log.Printf("Failed to update AI roast for PR #%d in %s: %v", body.GetNumber(), body.GetRepo().GetFullName(), err)
		return fmt.Errorf("Failed to update AI roast: %w", err)
	}
	log.Printf("Successfully closed PR #%d in %s", body.GetNumber(), body.GetRepo().GetFullName())
	return nil
}
