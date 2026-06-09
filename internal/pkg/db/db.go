package db

import (
	applog "chenze-faka/internal/pkg/logger"
	"fmt"
	"time"

	"chenze-faka/internal/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

var DB *gorm.DB

// IsReady 检查数据库是否真正可用（已连接且内部状态完整）
// 使用 panic recovery 防止 GORM 内部状态损坏时访问触发 panic
func IsReady() bool {
	if DB == nil {
		return false
	}
	// 捕获任何可能的 panic（GORM 内部状态损坏时 Where 等方法会 panic）
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

func Init(dsn string) error {
	var err error
	var dialector gorm.Dialector

	// 根据配置中的 driver 判断类型
	if config.AppConfig != nil && config.AppConfig.Database.IsSQLite() {
		dialector = sqlite.Open(dsn)
	} else if len(dsn) > 0 && (dsn[0] == '.' || dsn[0] == '/') {
		// 兜底：路径以 . 或 / 开头则视为 SQLite
		dialector = sqlite.Open(dsn)
	} else {
		// MySQL
		dialector = mysql.Open(dsn)
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

	if config.AppConfig != nil && config.AppConfig.Database.IsSQLite() {
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(5) // SQLite 允许少量并发连接用于读
	} else {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	DB = db
	applog.Infof("database connected")
	return nil
}
