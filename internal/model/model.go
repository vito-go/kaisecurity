package model

// ScanResultsInfo represents the results of a scan. primary key is ScanId
//
// CreateTime can be time.Time, or it will occur error:
// sql: Scan error on column index 12, name "create_time": unsupported Scan, storing driver.Value type string into type *time.Time;
// ScanMetadata and Summary are types of mtype.ScanMetadata and mtype.Summary, but they will be stored as strings in the database.
type ScanResultsInfo struct {
	ScanId       string `json:"scan_id"`
	Timestamp    string `json:"timestamp"`
	ScanStatus   string `json:"scan_status"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	Summary      string `json:"summary"`
	ScanMetadata string `json:"scan_metadata"`
	SourceFile   string `json:"source_file"`
	CreateTime   string `json:"create_time"`
	UpdateTime   string `json:"update_time"`
}

// Vulnerability represents a vulnerability found in a scan. primary key is a combination of ScanId and Id
// RiskFactors is a type of []string, but it will be stored as a string in the database.
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
	RiskFactors    string  `json:"risk_factors"`
	//
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}
