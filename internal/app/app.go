package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vito-go/kaisecurity/internal/dao"
	"github.com/vito-go/kaisecurity/internal/model"
	"github.com/vito-go/kaisecurity/pkg/db"
	"github.com/vito-go/mylog"
	"gorm.io/gorm"
	"io"
	"net/http"
	"sync"
	"time"
)

type APP struct {
	db     *gorm.DB
	allDao *dao.AllDao
}

func NewAppContext(dbPath string) (*APP, error) {
	gdb, err := db.NewSqliteDB(dbPath)
	if err != nil {
		return nil, err
	}
	err = gdb.Exec(model.SQLKaiSecurity).Error
	if err != nil {
		return nil, err
	}
	// 初始化表
	return &APP{
		db:     gdb,
		allDao: dao.NewAllDao(gdb),
	}, nil

}

type ScanRequest struct {
	Repo  string   `json:"repo"`
	Files []string `json:"files"`
}

func (a *APP) HandleScan(w http.ResponseWriter, r *http.Request) {
	//don't use r.Context() here, because when the request is canceled, the context will be canceled
	ctx := mylog.NewContext()
	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	const allowedRepo = "https://github.com/velancio/vulnerability_scans"
	if req.Repo != allowedRepo {
		http.Error(w, "unsupported repo", http.StatusBadRequest)
		return
	}
	if len(req.Files) == 0 {
		http.Error(w, "no files specified", http.StatusBadRequest)
		return
	}
	if len(req.Files) > 10 {
		http.Error(w, "too many files", http.StatusBadRequest)
		return
	}
	// 调用 service 扫描逻辑
	err := a.ProcessFilesConcurrently(ctx, req.Repo, req.Files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`success`))
}

func (a *APP) SaveScanData(ctx context.Context, items []ScanData) error {
	// update scanResultsInfo first
	for _, item := range items {
		m := item.ScanResults.ToScanResultsInfoModel()
		err := a.allDao.ScanResultsInfo.UpdateOrCreate(ctx, &m)
		if err != nil {
			return fmt.Errorf("failed to create scanResultsInfo: %w", err)
		}
		for _, vulnerability := range item.ScanResults.Vulnerabilities {
			vModel := vulnerability.ToVulnerabilityModel(item.ScanResults.ScanId)
			err = a.allDao.Vulnerability.UpdateOrCreate(ctx, &vModel)
			if err != nil {
				return fmt.Errorf("failed to create vulnerability: %w", err)
			}
		}
	}
	return nil
}
func (a *APP) ProcessFilesConcurrently(ctx context.Context, repo string, files []string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(files))
	sem := make(chan struct{}, 3) // 并发限制：最多 3 个文件同时处理
	for _, file := range files {
		wg.Add(1)
		sem <- struct{}{}
		go func(filename string) {
			defer wg.Done()
			defer func() { <-sem }()
			var err error
			defer func() {
				if err != nil {
					errCh <- err
				}
			}()
			data, err := DownloadFileWithRetry(repo, filename, 2)
			if err != nil {
				mylog.Ctx(ctx).Error(err)
				return
			}
			var result []ScanData
			err = json.Unmarshal(data, &result)
			if err != nil {
				mylog.Ctx(ctx).Error(err)
				return
			}
			if len(result) == 0 {
				mylog.Printf("no scan data found in %s", filename)
				return
			}
			err = a.SaveScanData(ctx, result)
			if err != nil {
				mylog.Ctx(ctx).Error(err)
				return
			}
		}(file)
	}
	wg.Wait()
	close(errCh)
	select {
	case err := <-errCh:
		// return an error if any of the goroutines failed
		return err
	default:
		return nil
	}
}

// DownloadFileWithRetry downloads a raw file from GitHub with retry logic
func DownloadFileWithRetry(repoURL, filename string, maxRetries int) ([]byte, error) {
	// 例：repoURL = https://github.com/velancio/vulnerability_scans
	// 原始文件 URL: https://raw.githubusercontent.com/velancio/vulnerability_scans/main/vulnscan1011.json
	rawURL, err := convertToRawGitHubURL(repoURL, "main", filename)
	if err != nil {
		return nil, err
	}
	for i := 0; i <= maxRetries; i++ {
		respBody, err := getRawURL(rawURL)
		if err != nil {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
			continue
		}
		return respBody, nil
	}
	return nil, fmt.Errorf("failed to download %s after %d attempts: %w", rawURL, maxRetries, err)
}
func getRawURL(rawURL string) ([]byte, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, status: %s", resp.StatusCode, resp.Status)
	}
	return io.ReadAll(resp.Body)
}
func convertToRawGitHubURL(repoURL, branch, filename string) (string, error) {
	const base = "https://raw.githubusercontent.com"
	var userRepo string
	_, err := fmt.Sscanf(repoURL, "https://github.com/%s", &userRepo)
	if err != nil {
		return "", fmt.Errorf("failed to parse repo URL: %w", err)
	}
	return fmt.Sprintf("%s/%s/%s/%s", base, userRepo, branch, filename), nil
}
