package handler

import (
	"io"
	"log"
	"net/http"

	middleware "github.com/chopstickleg/good-code/api/v1/_middleware"
	utils "github.com/chopstickleg/good-code/api/v1/_utils"
	handlers "github.com/chopstickleg/good-code/api/v1/_utils/handlers"

	"github.com/google/go-github/v72/github"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	middleware.AllowMethods(http.MethodPost)(GitHubWebhookHandler)(w, r)
}

func GitHubWebhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	if !utils.VerifyGitHubSignature(body, r.Header.Get("X-Hub-Signature-256")) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	eventType := github.WebHookType(r)
	if eventType == "" {
		log.Printf("Missing or empty X-GitHub-Event header")
		http.Error(w, "Missing event type", http.StatusBadRequest)
		return
	}

	log.Printf("Processing GitHub webhook event: %s", eventType)

	eventBody, err := github.ParseWebHook(eventType, body)
	if err != nil {
		log.Printf("Failed to parse GitHub event %s: %v", eventType, err)
		http.Error(w, "Unable to parse GitHub event", http.StatusBadRequest)
		return
	}
	switch eventType {
	case "pull_request":
		if prEvent, ok := eventBody.(*github.PullRequestEvent); ok {
			log.Printf("Processing pull request event for PR #%d", prEvent.GetNumber())
			handlers.HandlePullRequestEvent(w, *prEvent)
		} else {
			log.Printf("Failed to cast pull request event")
			http.Error(w, "Invalid pull request event", http.StatusBadRequest)
			return
		}
	case "installation":
		if installEvent, ok := eventBody.(*github.InstallationEvent); ok {
			log.Printf("Processing installation event: %s", installEvent.GetAction())
			handlers.HandleInstallationEvent(w, *installEvent)
		} else {
			log.Printf("Failed to cast installation event")
			http.Error(w, "Invalid installation event", http.StatusBadRequest)
			return
		}
	case "installation_target":
		if itEvent, ok := eventBody.(*github.InstallationTargetEvent); ok {
			log.Printf("Processing installation target event")
			handlers.HandleInstallationTargetEvent(w, *itEvent)
		} else {
			log.Printf("Failed to cast installation target event")
			http.Error(w, "Invalid installation target event", http.StatusBadRequest)
			return
		}
	case "repository":
		if repoEvent, ok := eventBody.(*github.RepositoryEvent); ok {
			log.Printf("Processing repository event: %s", repoEvent.GetAction())
			handlers.HandleRepositoryEvent(w, *repoEvent)
		} else {
			log.Printf("Failed to cast repository event")
			http.Error(w, "Invalid repository event", http.StatusBadRequest)
			return
		}
	default:
		log.Printf("Unsupported event type: %s", eventType)
		http.Error(w, "Unsupported event type", http.StatusBadRequest)
		return
	}

	log.Printf("Successfully processed %s event", eventType)
	w.Write([]byte("Event processed successfully"))
}
