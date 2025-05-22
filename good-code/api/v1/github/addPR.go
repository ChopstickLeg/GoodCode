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

	middleware "github.com/chopstickleg/good-code/api/v1/_middleware"

	"github.com/google/go-github/v72/github"
	"google.golang.org/genai"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	middleware.AllowMethods(http.MethodPost)(AddPRHandler)(w, r)
}

func AddPRHandler(w http.ResponseWriter, r *http.Request) {
	body := io.ReadAll(r.Body)
	if !VerifyGitHubSignature(body, r.Header.Get("X-Hub-Signature-256")) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}
	fmt.Println("Received valid request to add PR")
	token := os.Getenv("AI_API_TOKEN")
	if token == "" {
		http.Error(w, "Unable to get AI API token", http.StatusInternalServerError)
		return
	}

	var requestBody github.PullRequestEvent
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Unable to decode GitHub event", http.StatusBadRequest)
		return
	}

	// ghClient := github.NewClient(nil)

	// err := r.ParseMultipartForm(10 << 20)
	// if err != nil {
	// 	http.Error(w, "unable to parse form", http.StatusBadRequest)
	// 	return
	// }

	// file, _, err := r.FormFile("diff")
	// if err != nil {
	// 	http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
	// 	return
	// }
	// defer file.Close()

	// fileBytes, err := io.ReadAll(file)
	// if err != nil {
	// 	http.Error(w, "Unable to read file", http.StatusInternalServerError)
	// 	return
	// }

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
		genai.Text(string("placeholder text")),
		&config,
	)

	fmt.Print(result.Text())

	// reqBody := db.AI_request{
	// 	Prompt: string(fileBytes),
	// 	Model:  "gemma3:12b",
	// 	System: "You are a code review assistant. You will be given a diff of a pull request. Your task is to review the code and provide feedback. You should be sarcastic and condescending, but still helpful and provide useful feedback that is factually accurate to the best of your knowledge",
	// 	Stream: false,
	// }

	// body, err := json.Marshal(reqBody)
	// if err != nil {
	// 	http.Error(w, "Unable to marshal request", http.StatusInternalServerError)
	// 	return
	// }

	// req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	// if err != nil {
	// 	http.Error(w, "Unable to send request to AI API", http.StatusInternalServerError)
	// 	return
	// }
	// req.Header.Set("Authorization", "Bearer "+token)
	// req.Header.Set("Content-Type", "application/json")
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	http.Error(w, "Unable to send request to AI API", http.StatusInternalServerError)
	// 	return
	// }

	// defer resp.Body.Close()

	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "text/plain")
	// aiResponse, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	http.Error(w, "Unable to read response from AI API", http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(aiResponse)

	// gormdb, err := db.GetDB()
	// if err != nil {
	// 	fmt.Println("Failed to connect to database:", err)
	// 	return
	// }
	// pr := db.Pull_request{
	// 	ID:            0,
	// 	Author_id:     1,
	// 	Author_name:   "temp_author",
	// 	Source_branch: "temp_source",
	// 	Target_branch: "temp_target",
	// 	Has_comments:  true,
	// 	AIComments:    result.Text(),
	// }
	// gormdb.AutoMigrate(&db.Pull_request{})
	// gormdb.Create(&pr)
}

func VerifyGitHubSignature(payload []byte, signature string) bool {
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	key := hmac.New(sha256.New, []byte(secret))
	key.Write([]byte(string(payload)))
	computedSignature := fmt.Sprintf("sha256=%x", key.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(computedSignature))
}
