package db

import (
	applog "chenze-faka/internal/pkg/logger"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite" // pure-Go SQLite driver (registers as "sqlite")

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

var DB *gorm.DB
var driverType string // "mysql" 或 "sqlite"

func IsReady() bool {
	if DB == nil {
		return false
	}
	var ready bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				ready = false
			}
		}()
		_, err := DB.DB()
		ready = (err == nil)
	}()
	return ready
}

// Driver 返回当前数据库驱动类型 ("mysql" 或 "sqlite")
func Driver() string {
	return driverType
}

func Init(dsn string) error {
	var err error
	var dialector gorm.Dialector

	// 检测是否为 SQLite（.db 扩展名或不包含 MySQL 特定格式）
	isSQLite := strings.HasSuffix(dsn, ".db") || strings.HasSuffix(dsn, ".sqlite") ||
		(len(dsn) < 30 && strings.HasPrefix(dsn, "./") && !strings.Contains(dsn, "@tcp"))

	if isSQLite {
		// SQLite：使用 modernc.org/sqlite（纯 Go，无 CGO）
		// 通过 DriverName=sqlite 让 gorm 走 modernc 的驱动注册
		dialector = sqlite.Dialector{
			DSN:        dsn,
			DriverName: "sqlite",
		}
		driverType = "sqlite"
	} else {
		// MySQL
		dialector = mysql.Open(dsn)
		driverType = "mysql"
	}

	var db *gorm.DB
	db, err = gorm.Open(dialector, &gorm.Config{
		Logger: gormlog.Default.LogMode(gormlog.Info),
	})
	if err != nil {
		DB = nil
		return fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil || sqlDB == nil {
		DB = nil
		return fmt.Errorf("get sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	applog.Infof("database connected")
	return nil
}
