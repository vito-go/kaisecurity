package httpsrv

import (
	"encoding/json"
	"github.com/vito-go/mylog"
	"net/http"
)

type QueryRequest struct {
	Filters filters `json:"filters"`
}
type filters struct {
	Severity string `json:"severity"`
}

func (s *Server) HandleQuery(w http.ResponseWriter, r *http.Request) {
	ctx := mylog.NewContext()
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Filters.Severity == "" {
		http.Error(w, "no severity specified", http.StatusBadRequest)
		return
	}
	// do something
	items, err := s.allDao.Vulnerability.ItemsBySeverity(ctx, req.Filters.Severity)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, "encode response failed", http.StatusInternalServerError)
		return
	}
}
