package testutil

import (
	"context"
	"encoding/json"
	"github.com/vito-go/kaisecurity/internal/dao"
	"github.com/vito-go/kaisecurity/internal/httpsrv"
	"github.com/vito-go/kaisecurity/internal/model"
	"github.com/vito-go/kaisecurity/pkg/db"
	"testing"
)

func NewServer(t *testing.T) *httpsrv.Server {
	t.Helper()
	ap, err := httpsrv.NewServer(":memory:")
	if err != nil {
		t.Fatalf("failed to init app: %v", err)
	}
	return ap
}

func NewAllDao(t *testing.T) *dao.AllDao {
	gdb, err := db.NewSqliteDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	err = gdb.Exec(model.SQLKaiSecurity).Error
	if err != nil {
		t.Fatal(err)
	}
	return dao.NewAllDao(gdb)
}
func InsertTestVulnerability(t *testing.T, dao *dao.AllDao) model.Vulnerability {
	t.Helper()
	ctx := context.Background()
	data := `{
		"scan_id": "scan_id_1",
		"id": "CVE-2024-2222",
		"severity": "HIGH",
		"cvss": 8.2,
		"status": "active",
		"package_name": "spring-security",
		"current_version": "5.6.0",
		"fixed_version": "5.6.1",
		"description": "Authentication bypass in Spring Security",
		"published_date": "2025-01-27T00:00:00Z",
		"link": "https://nvd.nist.gov/vuln/detail/CVE-2024-2222",
		"risk_factors": "['CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H']"
	}`

	var vuln model.Vulnerability
	if err := json.Unmarshal([]byte(data), &vuln); err != nil {
		t.Fatalf("failed to unmarshal test vuln: %v", err)
	}
	if err := dao.Vulnerability.UpdateOrCreate(ctx, &vuln); err != nil {
		t.Fatalf("failed to insert test vuln: %v", err)
	}
	return vuln
}
