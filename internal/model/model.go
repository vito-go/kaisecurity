package model

// ScanResultsInfo represents the results of a scan. primary key is ScanId
//
//	CreateTime can be time.Time sql: Scan error on column index 12, name "create_time": unsupported Scan, storing driver.Value type string into type *time.Time;
type ScanResultsInfo struct {
	ScanId       string `json:"scan_id"`
	Timestamp    string `json:"timestamp"`
	ScanStatus   string `json:"scan_status"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	Summary      string `json:"summary"`       // mtype.Summary
	ScanMetadata string `json:"scan_metadata"` // mtype.ScanMetadata
	CreateTime   string `json:"create_time"`
	UpdateTime   string `json:"update_time"`
}

// Vulnerability represents a vulnerability found in a scan. primary key is a combination of ScanId and Id
type Vulnerability struct {
	ScanId         string  `json:"scan_id"` // reference to ScanResultsInfo.ScanId
	Id             string  `json:"id"`
	Severity       string  `json:"severity"`
	Cvss           float64 `json:"cvss"`
	Status         string  `json:"status"`
	PackageName    string  `json:"package_name"`
	CurrentVersion string  `json:"current_version"`
	FixedVersion   string  `json:"fixed_version"`
	Description    string  `json:"description"`
	PublishedDate  string  `json:"published_date"`
	Link           string  `json:"link"`
	RiskFactors    string  `json:"risk_factors"` //  []string
	//
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}
