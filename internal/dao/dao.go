package dao

import "gorm.io/gorm"

type AllDao struct {
	ScanResultsInfo *scanResultsInfo
	Vulnerability   *vulnerability
}

func NewAllDao(db *gorm.DB) *AllDao {
	return &AllDao{
		ScanResultsInfo: &scanResultsInfo{Gdb: db},
		Vulnerability:   &vulnerability{Gdb: db},
	}
}
