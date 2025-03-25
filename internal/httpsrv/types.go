package httpsrv

import (
	"encoding/json"
	"github.com/vito-go/kaisecurity/internal/model"
	"github.com/vito-go/kaisecurity/internal/model/mtype.go"
	"time"
)

type ScanData struct {
	ScanResults ScanResults `json:"scanResults"`
}

type Vulnerability struct {
	Id             string   `json:"id"`
	Severity       string   `json:"severity"`
	Cvss           float64  `json:"cvss"`
	Status         string   `json:"status"`
	PackageName    string   `json:"package_name"`
	CurrentVersion string   `json:"current_version"`
	FixedVersion   string   `json:"fixed_version"`
	Description    string   `json:"description"`
	PublishedDate  string   `json:"published_date"`
	Link           string   `json:"link"`
	RiskFactors    []string `json:"risk_factors"`
}

func (v *Vulnerability) ToVulnerabilityModel(scanId string) model.Vulnerability {
	b, _ := json.Marshal(v.RiskFactors)
	return model.Vulnerability{
		ScanId:         scanId,
		Id:             v.Id,
		Severity:       v.Severity,
		Cvss:           v.Cvss,
		Status:         v.Status,
		PackageName:    v.PackageName,
		CurrentVersion: v.CurrentVersion,
		FixedVersion:   v.FixedVersion,
		Description:    v.Description,
		PublishedDate:  v.PublishedDate,
		Link:           v.Link,
		RiskFactors:    string(b),
		CreateTime:     time.Now().Format(time.RFC3339),
		UpdateTime:     time.Now().Format(time.RFC3339),
	}
}

type ScanResults struct {
	ScanId          string             `json:"scan_id"`
	Timestamp       string             `json:"timestamp"`
	ScanStatus      string             `json:"scan_status"`
	ResourceType    string             `json:"resource_type"`
	ResourceName    string             `json:"resource_name"`
	Vulnerabilities []Vulnerability    `json:"vulnerabilities"`
	Summary         mtype.Summary      `json:"summary"`
	ScanMetadata    mtype.ScanMetadata `json:"scan_metadata"`
}

func (s *ScanResults) ToScanResultsInfoModel(sourceFile string) model.ScanResultsInfo {
	return model.ScanResultsInfo{
		ScanId:       s.ScanId,
		Timestamp:    s.Timestamp,
		ScanStatus:   s.ScanStatus,
		ResourceType: s.ResourceType,
		ResourceName: s.ResourceName,
		SourceFile:   sourceFile,
		Summary:      s.Summary.Marshal(),
		ScanMetadata: s.ScanMetadata.Marshal(),
		CreateTime:   time.Now().Format(time.RFC3339),
		UpdateTime:   time.Now().Format(time.RFC3339),
	}
}
