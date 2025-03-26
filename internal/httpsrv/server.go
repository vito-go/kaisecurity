package httpsrv

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
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	db      *gorm.DB
	allDao  *dao.AllDao
	httpSrv *http.Server
}

func (s *Server) AllDao() *dao.AllDao {
	return s.allDao
}

func NewServer(dbPath string) (*Server, error) {
	gdb, err := db.NewSqliteDB(dbPath)
	if err != nil {
		return nil, err
	}
	err = gdb.Exec(model.SQLKaiSecurity).Error
	if err != nil {
		return nil, err
	}
	// 初始化表
	return &Server{
		db:     gdb,
		allDao: dao.NewAllDao(gdb),
	}, nil

}

func (s *Server) SaveScanData(ctx context.Context, sourceFile string, items []ScanData) error {
	// update scanResultsInfo first
	for _, item := range items {
		m := item.ScanResults.ToScanResultsInfoModel(sourceFile)
		err := s.allDao.ScanResultsInfo.UpdateOrCreate(ctx, &m)
		if err != nil {
			return fmt.Errorf("failed to create scanResultsInfo: %w", err)
		}
		for _, vulnerability := range item.ScanResults.Vulnerabilities {
			vModel := vulnerability.ToVulnerabilityModel(item.ScanResults.ScanId)
			err = s.allDao.Vulnerability.UpdateOrCreate(ctx, &vModel)
			if err != nil {
				return fmt.Errorf("failed to create vulnerability: %w", err)
			}
		}
	}
	return nil
}
func (s *Server) ProcessFilesConcurrently(ctx context.Context, repo string, files []string, getFile func(repoURL, branch, filename string, maxRetries int) (io.ReadCloser, error)) error {
	const branch = "main"
	var wg sync.WaitGroup
	errCh := make(chan error, len(files))
	sem := make(chan struct{}, 3) // control the number of concurrent goroutines
	for _, file := range files {
		wg.Add(1)
		sem <- struct{}{}
		go func(sourceFile string) {
			defer wg.Done()
			defer func() { <-sem }()
			var err error
			defer func() {
				if err != nil {
					errCh <- err
				}
			}()
			body, err := getFile(repo, branch, sourceFile, 2)
			if err != nil {
				mylog.Ctx(ctx).Error(err)
				return
			}
			var result []ScanData
			err = json.NewDecoder(body).Decode(&result)
			if err != nil {
				mylog.Ctx(ctx).Error(err)
				return
			}
			if len(result) == 0 {
				mylog.Printf("no scan data found in %s", sourceFile)
				return
			}
			err = s.SaveScanData(ctx, sourceFile, result)
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

func HandleFunc(mux *http.ServeMux, pattern string, handler http.HandlerFunc) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &responseWriter{ResponseWriter: w, statusCode: 200}
		defer func() {
			latency := time.Since(start)
			remoteIP := r.RemoteAddr
			// record the status code, non-200 status could be recorded in a separate log file if needed.
			//if lrw.statusCode != http.StatusOK {

			//}
			mylog.Ctx(r.Context()).Infof("%s -> %s %d %s %s %s %s %s",
				remoteIP,
				r.Host,
				lrw.statusCode,
				r.Method,
				r.URL.Path,
				r.URL.RawQuery,
				formatBytes(lrw.bytesWritten),
				latency,
			)
		}()
		handler(lrw, r)
	})
}
func (s *Server) StartServer(port uint) error {
	mux := http.NewServeMux()
	HandleFunc(mux, "POST /scan", s.HandleScan)
	HandleFunc(mux, "POST /query", s.HandleQuery)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	s.httpSrv = server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	// real port
	_, realPort, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		return err
	}
	mylog.Printf("Server is running on :%s", realPort)
	return s.httpSrv.Serve(lis)
}

// ShutDownServer gracefully shuts down the server
func (s *Server) ShutDownServer(ctx context.Context) error {
	if s.httpSrv == nil {
		return nil
	}
	return s.httpSrv.Shutdown(ctx)
}
