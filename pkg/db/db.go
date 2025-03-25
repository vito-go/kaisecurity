package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSqliteDB(dbPath string) (*gorm.DB, error) {
	GDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		CreateBatchSize: 500,
		Logger:          logger.Default.LogMode(logger.Warn), // if needed,    use logger.Info can be used
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
