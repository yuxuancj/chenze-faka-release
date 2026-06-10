package db

import (
	applog "chenze-faka/internal/pkg/logger"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

var DB *gorm.DB

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

func Init(dsn string) error {
	var err error
	var dialector gorm.Dialector

	dialector = mysql.Open(dsn)

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
