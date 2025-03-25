package dao

import (
	"context"
	"github.com/vito-go/kaisecurity/internal/model"
	"gorm.io/gorm"
)

type scanResultsInfo struct {
	Gdb *gorm.DB
}

func (s *scanResultsInfo) Table() string {
	return "scan_results_info"
}

// ItemByScanId get scan results info by scan id
func (s *scanResultsInfo) ItemByScanId(ctx context.Context, scanId string) (*model.ScanResultsInfo, error) {
	var m model.ScanResultsInfo
	tx := s.Gdb.WithContext(ctx).Table(s.Table()).Where("scan_id = ?", scanId).First(&m)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &m, nil
}
func (s *scanResultsInfo) UpdateOrCreate(ctx context.Context, m *model.ScanResultsInfo) (err error) {
	TX := s.Gdb.WithContext(ctx).Table(s.Table()).Begin()
	defer func() {
		if err != nil {
			TX.Rollback()
			return
		}
		err = TX.Commit().Error
	}()
	tx := TX.Select("*").Omit("create_time").Where("scan_id = ?", m.ScanId).Updates(m)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		if err = TX.Create(m).Error; err != nil {
			return err
		}
		return nil
	}
	return nil
}
