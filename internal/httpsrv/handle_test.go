package httpsrv_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/vito-go/kaisecurity/internal/httpsrv"
	"github.com/vito-go/kaisecurity/internal/testutil"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func performRequest(t *testing.T, handlerFunc http.HandlerFunc, payload string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest("POST", "/scan", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handlerFunc(w, req)
	return w
}

func testNewApp(t *testing.T) *httpsrv.Server {
	t.Helper()
	ap, err := httpsrv.NewServer(":memory:")
	if err != nil {
		t.Fatalf("failed to init app: %v", err)
	}
	return ap
}

func TestHandleScan(t *testing.T) {
	app := testNewApp(t)

	t.Run("success", func(t *testing.T) {
		payload := `{
			"repo": "https://github.com/velancio/vulnerability_scans",
			"files": ["vulnscan15.json"]
		}`
		rr := performRequest(t, app.HandleScan, payload)
		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rr.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		rr := performRequest(t, app.HandleScan, "not-json")
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 BadRequest for invalid JSON, got %d", rr.Code)
		}
	})

	t.Run("no files specified", func(t *testing.T) {
		payload := `{
			"repo": "https://github.com/velancio/vulnerability_scans",
			"files": []
		}`
		rr := performRequest(t, app.HandleScan, payload)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 BadRequest for empty files, got %d", rr.Code)
		}
	})

	t.Run("unsupported repo", func(t *testing.T) {
		payload := `{
			"repo": "https://github.com/unknown/other_repo",
			"files": ["somefile.json"]
		}`
		rr := performRequest(t, app.HandleScan, payload)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 BadRequest for unsupported repo, got %d", rr.Code)
		}
	})
}

func TestHandleQuery(t *testing.T) {
	newApp := testNewApp(t)

	t.Run("valid data", func(t *testing.T) {
		// insert some data, so we can query it
		testutil.InsertTestVulnerability(t, newApp.AllDao())
	})

	t.Run("missing severity", func(t *testing.T) {
		payload := `{"filters": {}}`
		res := performRequest(t, newApp.HandleQuery, payload)
		if res.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for missing severity, got %d", res.Code)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		res := performRequest(t, newApp.HandleQuery, `not-json`)
		if res.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for invalid json, got %d", res.Code)
		}
	})
	t.Run("success", func(t *testing.T) {
		payload := `{"filters": {"severity": "HIGH"}}`
		res := performRequest(t, newApp.HandleQuery, payload)
		if res.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", res.Code)
		}
		reuslt, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(string(reuslt), "HIGH") {
			t.Fatalf("expected response to contain 'HIGH', got: %s", string(reuslt))
		}
	})
}
func TestServer_saveScanData(t *testing.T) {
	srv := testutil.NewServer(t)
	ctx := context.Background()

	scan := httpsrv.ScanData{
		ScanResults: httpsrv.ScanResults{
			ScanId:       "scan_test_1",
			Timestamp:    time.Now().Format(time.RFC3339),
			ScanStatus:   "completed",
			ResourceType: "container",
			ResourceName: "example:1.0",
			Vulnerabilities: []httpsrv.Vulnerability{
				{
					Id:          "CVE-123",
					Severity:    "HIGH",
					Cvss:        8.2,
					PackageName: "openssl",
					Status:      "active",
				},
			},
		},
	}
	err := srv.SaveScanData(ctx, "test.json", []httpsrv.ScanData{scan})
	if err != nil {
		t.Fatalf("saveScanData failed: %v", err)
	}
	vulns, _ := srv.AllDao().Vulnerability.ItemsBySeverity(ctx, "HIGH")

	if len(vulns) != 1 || vulns[0].Id != "CVE-123" {
		t.Errorf("expected vulnerability to be saved, got %+v", vulns)
	}
}

func TestServer_processFilesConcurrently(t *testing.T) {
	// mock GetGithubFile 函数
	getMockFile := func(repo, branch, filename string, maxRetries int) (io.ReadCloser, error) {
		scanData := []httpsrv.ScanData{
			{
				ScanResults: httpsrv.ScanResults{
					ScanId:       "scan_test_2",
					Timestamp:    time.Now().Format(time.RFC3339),
					ResourceType: "container",
					ResourceName: "mocked",
					Vulnerabilities: []httpsrv.Vulnerability{
						{Id: "CVE-999", Severity: "HIGH", Cvss: 7.9, Status: "active"},
					},
				},
			},
		}
		bs, _ := json.Marshal(scanData)
		return io.NopCloser(bytes.NewReader(bs)), nil
	}

	srv := testutil.NewServer(t)

	err := srv.ProcessFilesConcurrently(context.Background(),
		"https://mock.com/repo", []string{"fake1.json", "fake2.json"}, getMockFile)
	if err != nil {
		t.Fatalf("processFilesConcurrently failed: %v", err)
	}

	// 验证是否写入了数据
	vulns, _ := srv.AllDao().Vulnerability.ItemsBySeverity(context.Background(), "HIGH")
	if len(vulns) == 0 {
		t.Error("expected vulnerability inserted, got none")
	}
	var found bool
	for _, v := range vulns {
		if v.Id == "CVE-999" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected vulnerability with ID CVE-999, got %+v", vulns)
	}
	found = false
	//Verify if the scan results are written
	scans, err := srv.AllDao().ScanResultsInfo.ItemByScanId(context.Background(), "scan_test_2")
	if err != nil {
		t.Fatalf("failed to query scan results info: %v", err)
	}
	if scans.ScanId != "scan_test_2" {
		t.Errorf("expected scan result with ID scan_test_2, got %+v", scans)
	}
}
