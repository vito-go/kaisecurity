package model

import (
	"time"
)

// ScanResultsInfo represents the results of a scan. primary key is ScanId
type ScanResultsInfo struct {
	ScanId       string    `json:"scan_id"`
	Timestamp    time.Time `json:"timestamp"`
	ScanStatus   string    `json:"scan_status"`
	ResourceType string    `json:"resource_type"`
	ResourceName string    `json:"resource_name"`
	Summary      string    `json:"summary"`       // mtype.Summary
	ScanMetadata string    `json:"scan_metadata"` // mtype.ScanMetadata
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
}

// Vulnerability represents a vulnerability found in a scan. primary key is a combination of ScanId and Id
type Vulnerability struct {
	ScanId         string    `json:"scan_id"` // reference to ScanResultsInfo.ScanId
	Id             string    `json:"id"`
	Severity       string    `json:"severity"`
	Cvss           float64   `json:"cvss"`
	Status         string    `json:"status"`
	PackageName    string    `json:"package_name"`
	CurrentVersion string    `json:"current_version"`
	FixedVersion   string    `json:"fixed_version"`
	Description    string    `json:"description"`
	PublishedDate  time.Time `json:"published_date"`
	Link           string    `json:"link"`
	RiskFactors    string    `json:"risk_factors"` //  []string
	//
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
