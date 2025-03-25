package util

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// GetGithubFile downloads a raw file from GitHub with retry logic
// remember to close the response body after using it
func GetGithubFile(repoURL, branch, filename string, maxRetries int) (io.ReadCloser, error) {
	// repoURL = https://github.com/velancio/vulnerability_scans
	// File URL: https://raw.githubusercontent.com/velancio/vulnerability_scans/main/vulnscan1011.json
	rawURL, err := ConvertToRawGitHubURL(repoURL, branch, filename)
	if err != nil {
		return nil, err
	}
	for i := 0; i <= maxRetries; i++ {
		resp, err := http.Get(rawURL)
		if err != nil {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code: %d, status: %s", resp.StatusCode, resp.Status)
		}
		return resp.Body, nil
	}
	return nil, fmt.Errorf("failed to download %s after %d attempts: %w", rawURL, maxRetries, err)
}

func ConvertToRawGitHubURL(repoURL, branch, filename string) (string, error) {
	const base = "https://raw.githubusercontent.com"
	var userRepo string
	_, err := fmt.Sscanf(repoURL, "https://github.com/%s", &userRepo)
	if err != nil {
		return "", fmt.Errorf("failed to parse repo URL: %w", err)
	}
	return fmt.Sprintf("%s/%s/%s/%s", base, userRepo, branch, filename), nil
}
