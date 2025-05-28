package handler

import (
	"io"
	"net/http"

	middleware "github.com/chopstickleg/good-code/api/v1/_middleware"
	utils "github.com/chopstickleg/good-code/api/v1/_utils"

	"github.com/google/go-github/v72/github"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	middleware.AllowMethods(http.MethodPost)(GitHubWebhookHandler)(w, r)
}

func GitHubWebhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	if !utils.VerifyGitHubSignature(body, r.Header.Get("X-Hub-Signature-256")) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}
	eventType := github.WebHookType(r)

	eventBody, err := github.ParseWebHook(eventType, body)
	if err != nil {
		http.Error(w, "Unable to parse GitHub event", http.StatusBadRequest)
		return
	}
	switch eventType {
	case "pull_request":
		if prEvent, ok := eventBody.(*github.PullRequestEvent); ok {
			utils.AddPRHandler(w, *prEvent)
		} else {
			http.Error(w, "Invalid pull request event", http.StatusBadRequest)
			return
		}
	case "installation_target":
		//handle that
	case "meta":
		//handle that
	default:
		http.Error(w, "Unsupported event type", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Event processed successfully"))
}