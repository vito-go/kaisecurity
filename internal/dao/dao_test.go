package dao_test

import (
	"context"
	"fmt"
	"github.com/vito-go/kaisecurity/internal/model"
	"github.com/vito-go/kaisecurity/internal/testutil"
	"testing"
	"time"
)

var ctx = context.Background()

func Test_scanResultsInfo_dao(t *testing.T) {
	_allDao := testutil.NewAllDao(t)
	m := &model.ScanResultsInfo{
		ScanId:       fmt.Sprintf("scan_id_%d", 1),
		Timestamp:    "2025-03-25T09:15:00Z",
		ScanStatus:   "completed",
		ResourceType: "container",
		ResourceName: "auth-service:2.1.0",
		Summary:      "",
		ScanMetadata: "",
		SourceFile:   "",
		CreateTime:   time.Now().Format(time.RFC3339),
		UpdateTime:   time.Now().Format(time.RFC3339),
	}
	err := _allDao.ScanResultsInfo.UpdateOrCreate(ctx, m)
	if err != nil {
		t.Error(err)
	}
	result, err := _allDao.ScanResultsInfo.ItemByScanId(ctx, m.ScanId)
	if err != nil {
		t.Error(err)
	}
	if result.ScanId != m.ScanId {
		t.Errorf("expected %s, got %s", m.ScanId, result.ScanId)
	}
}

func Test_vulnerability_dao(t *testing.T) {
	_allDao := testutil.NewAllDao(t)
	m := testutil.InsertTestVulnerability(t, _allDao)
	items, err := _allDao.Vulnerability.ItemsBySeverity(ctx, m.Severity)
	if err != nil {
		t.Error(err)
	}
	if len(items) == 0 {
		t.Error("no data found")
	}
	var found bool
	for _, item := range items {
		if item.ScanId == m.ScanId && item.Id == m.Id {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("vulnerability not found in result: %+v", m)
	}
}
