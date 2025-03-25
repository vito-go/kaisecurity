package dao

import (
	"context"
	"github.com/vito-go/kaisecurity/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type scanResultsInfo struct {
	Gdb *gorm.DB
}

func (s *scanResultsInfo) Table() string {
	return "scan_results_info"
}

// CreateBatch .
func (s *scanResultsInfo) CreateBatch(ctx context.Context, items ...model.ScanResultsInfo) error {
	if len(items) == 0 {
		return nil
	}
	return s.Gdb.WithContext(ctx).Table(s.Table()).Clauses(clause.Insert{
		Modifier: "OR IGNORE",
	}).Create(items).Error
}

func (s *scanResultsInfo) ItemByScanId(ctx context.Context, scanId string) (*model.ScanResultsInfo, error) {
	var msg model.ScanResultsInfo
	tx := s.Gdb.WithContext(ctx).Table(s.Table()).Where("scan_id = ?", scanId).First(&msg)
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &msg, tx.Error
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
	tx := TX.Where("scan_id = ?", m.ScanId).Updates(m)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		if err = tx.Create(m).Error; err != nil {
			return err
		}
		return nil
	}
	return nil
}
