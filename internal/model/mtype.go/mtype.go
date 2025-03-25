package mtype

import "encoding/json"

type SeverityCounts struct {
	CRITICAL int `json:"CRITICAL"`
	HIGH     int `json:"HIGH"`
	MEDIUM   int `json:"MEDIUM"`
	LOW      int `json:"LOW"`
}
type Summary struct {
	TotalVulnerabilities int            `json:"total_vulnerabilities"`
	SeverityCounts       SeverityCounts `json:"severity_counts"`
	FixableCount         int            `json:"fixable_count"`
	Compliant            bool           `json:"compliant"`
}

func (s Summary) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

type ScanMetadata struct {
	ScannerVersion  string   `json:"scanner_version"`
	PoliciesVersion string   `json:"policies_version"`
	ScanningRules   []string `json:"scanning_rules"`
	ExcludedPaths   []string `json:"excluded_paths"`
}

func (s ScanMetadata) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}
