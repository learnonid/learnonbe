package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"learnonbe/model"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

// UploadToGithub uploads a file to GitHub repository
func UploadToGithub(fileName, content string) error {
	// Ambil owner dan repo dari environment
	var githubOwner = os.Getenv("GITHUB_OWNER")
	var githubRepo = os.Getenv("GITHUB_REPO")
	var githubToken = os.Getenv("GITHUB_TOKEN")

	// Validasi environment variables
	if githubOwner == "" || githubRepo == "" || githubToken == "" {
		return fmt.Errorf("GITHUB_OWNER, GITHUB_REPO, or GITHUB_TOKEN is not set in the environment")
	}

	// Buat URL API
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", githubOwner, githubRepo, fileName)

	// Buat payload permintaan
	uploadRequest := model.GithubUploadRequest{
		Message: "Upload payment receipt " + fileName,
		Content: content,
	}

	jsonData, err := json.Marshal(uploadRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal upload request: %v", err)
	}

	// Buat HTTP request
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Content-Type", "application/json")

	// Kirim permintaan
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Periksa status code
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload file to Github: %s, response: %s", resp.Status, string(body))
	}

	return nil
}

// get all user event registration repository
func GetAllUERegistration(ctx context.Context, db *mongo.Database) ([]model.UserEventRegistration, error) {
	collection := db.Collection("ueregist")
	cursor, err := collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var registrations []model.UserEventRegistration
	if err = cursor.All(ctx, &registrations); err != nil {
		return nil, err
	}

	return registrations, nil
}