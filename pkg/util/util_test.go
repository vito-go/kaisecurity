package util_test

import (
	"bytes"
	"github.com/jarcoal/httpmock"
	"github.com/vito-go/kaisecurity/pkg/util"
	"io"
	"net/http"
	"testing"
)

const (
	repoURL  = "https://github.com/velancio/vulnerability_scans"
	branch   = "main"
	filename = "vulnscan1011.json"
)

var maxRetries = 2

func TestConvertToRawGitHubURL(t *testing.T) {
	expected := "https://raw.githubusercontent.com/velancio/vulnerability_scans/main/vulnscan1011.json"
	url, err := util.ConvertToRawGitHubURL(repoURL, branch, filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

func TestConvertToRawGitHubURL_Invalid(t *testing.T) {
	_, err := util.ConvertToRawGitHubURL("invalid-url", branch, filename)
	if err == nil {
		t.Error("expected error for invalid repo URL, got nil")
	}
}

func TestGetGithubFile_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	rawURL, _ := util.ConvertToRawGitHubURL(repoURL, branch, filename)
	httpmock.RegisterResponder("GET", rawURL,
		httpmock.NewStringResponder(200, `{"status":"ok"}`))

	body, err := util.GetGithubFile(repoURL, branch, filename, maxRetries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer body.Close()

	data, _ := io.ReadAll(body)
	if !bytes.Contains(data, []byte("ok")) {
		t.Errorf("expected response body to contain 'ok', got: %s", string(data))
	}
}

func TestGetGithubFile_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	rawURL, _ := util.ConvertToRawGitHubURL(repoURL, branch, filename)
	httpmock.RegisterResponder("GET", rawURL,
		httpmock.NewStringResponder(404, "Not Found"))

	_, err := util.GetGithubFile(repoURL, branch, filename, maxRetries)
	if err == nil {
		t.Error("expected error for 404 response, got nil")
	}
}

func TestGetGithubFile_RetrySuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	rawURL, _ := util.ConvertToRawGitHubURL(repoURL, branch, filename)
	callCount := 0

	httpmock.RegisterResponder("GET", rawURL, func(req *http.Request) (*http.Response, error) {
		callCount++
		if callCount < maxRetries {
			return nil, io.ErrUnexpectedEOF
		}
		return httpmock.NewStringResponse(200, `{"status":"retry-success"}`), nil
	})

	body, err := util.GetGithubFile(repoURL, branch, filename, maxRetries)
	if err != nil {
		t.Fatalf("expected success on retry, got error: %v", err)
	}
	defer body.Close()

	data, _ := io.ReadAll(body)
	if !bytes.Contains(data, []byte("retry-success")) {
		t.Errorf("expected retry success response, got: %s", string(data))
	}
	if callCount != maxRetries {
		t.Errorf("expected %d calls, got %d", maxRetries, callCount)
	}
}
