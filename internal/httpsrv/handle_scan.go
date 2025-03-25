package httpsrv

import (
	"encoding/json"
	"github.com/vito-go/kaisecurity/pkg/util"
	"github.com/vito-go/mylog"
	"net/http"
)

type ScanRequest struct {
	Repo  string   `json:"repo"`
	Files []string `json:"files"`
}

func (s *Server) HandleScan(w http.ResponseWriter, r *http.Request) {
	//don't use r.Context() here, because when the request is canceled, the context will be canceled,
	// and the subsequent processing will be affected
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
	err := s.ProcessFilesConcurrently(ctx, req.Repo, req.Files, util.GetGithubFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`success`))
}
