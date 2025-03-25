package app

import "testing"

func Test_convertToRawGitHubURL(t *testing.T) {
	repoURL := "https://github.com/velancio/vulnerability_scans"
	result, err := convertToRawGitHubURL(repoURL, "main", "vulnscan1011.json")
	if err != nil {
		t.Fatal(err)
	}

	// expected: https://raw.githubusercontent.com/velancio/vulnerability_scans/main/vulnscan1011.json
	if result != "https://raw.githubusercontent.com/velancio/vulnerability_scans/main/vulnscan1011.json" {
		t.Fatalf("unexpected result: %s", result)
	}

}
