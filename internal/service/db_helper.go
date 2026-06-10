package service

import (
	"chenze-faka/internal/pkg/db"

	"gorm.io/gorm"
)

// GetDB 导出数据库连接（用于高级功能的仪表盘查询）
func GetDB() *gorm.DB {
	return db.DB
}
