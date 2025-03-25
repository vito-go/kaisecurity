package db

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSqliteDB(dbPath string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s", dbPath)
	GDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		CreateBatchSize: 500,
		Logger:          logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db, err := GDB.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return GDB, nil
}
