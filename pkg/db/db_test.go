package db_test

import (
	"os"
	"testing"

	"github.com/vito-go/kaisecurity/pkg/db"
)

func TestNewSqliteDBFile(t *testing.T) {
	// 使用内存数据库
	gdb, err := db.NewSqliteDB("/&/！not db") // use invalid path
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	dbPath := ".test.db"
	defer os.Remove(dbPath)
	gdb, err = db.NewSqliteDB(dbPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 执行一个简单的 SQL 来验证连接
	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("database ping failed: %v", err)
	}
	// 关闭数据库
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("failed to close database: %v", err)
	}
}

func TestNewSqliteDBMemory(t *testing.T) {
	gdb, err := db.NewSqliteDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gdb == nil {
		t.Fatal("expected non-nil *gorm.DB")
	}

	// 执行一个简单的 SQL 来验证连接
	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("database ping failed: %v", err)
	}
	// 关闭数据库
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("failed to close database: %v", err)
	}
}
