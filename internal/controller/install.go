package controller

import (
	appcfg "chenze-faka/config"
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/config"
	"chenze-faka/internal/pkg/db"
	"chenze-faka/internal/pkg/response"
	"chenze-faka/internal/service"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type InstallController struct{}

func NewInstallController() *InstallController { return &InstallController{} }

// EnvCheck 环境检测
func (ic *InstallController) EnvCheck(ctx *gin.Context) {
	// 检测目录写权限
	writeable := true
	if err := os.WriteFile("install_test.tmp", []byte("test"), 0644); err == nil {
		os.Remove("install_test.tmp")
	} else {
		writeable = false
	}

	// 数据库检测
	dbStatus := "✅"
	if db.DB != nil {
		sqlDB, _ := db.DB.DB()
		if sqlDB != nil {
			if err := sqlDB.Ping(); err != nil {
				dbStatus = "❌ " + err.Error()
			}
		}
	} else {
		dbStatus = "⚠️ 未连接（安装时将自动连接）"
	}

	response.Success(ctx, gin.H{
		"go_version": runtime.Version(),
		"writeable":  writeable,
		"db_status":  dbStatus,
	})
}

// Install 执行安装
func (ic *InstallController) Install(ctx *gin.Context) {
	var req struct {
		Driver     string `json:"driver"`      // mysql 或 sqlite，默认 mysql
		DBHost     string `json:"db_host"`    // MySQL 主机
		DBPort     int    `json:"db_port"`    // MySQL 端口
		DBUser     string `json:"db_user"`    // MySQL 用户名
		DBPassword string `json:"db_password"` // MySQL 密码
		DBName     string `json:"db_name"`    // 数据库名（MySQL）或文件名（SQLite）
		AdminUser  string `json:"admin_user"`
		AdminPass  string `json:"admin_pass"`
		AdminEmail string `json:"admin_email"`
		LicenseKey string `json:"license_key"` // 授权码
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, response.CodeParamError, "参数错误: "+err.Error())
		return
	}

	// 参数默认值
	if req.Driver == "" {
		req.Driver = "mysql" // 默认 MySQL
	}
	if req.DBHost == "" {
		req.DBHost = "127.0.0.1"
	}
	if req.DBPort == 0 {
		req.DBPort = 3306
	}
	if req.DBName == "" {
		req.DBName = "chenze_faka"
	}
	if req.AdminUser == "" || req.AdminPass == "" || req.AdminEmail == "" {
		response.Error(ctx, response.CodeParamError, "管理员信息不完整")
		return
	}

	// 授权码验证（调用授权站）
	licenseSvc := service.NewLicenseService()
	siteURL := licenseSvc.DetectURL(ctx.Request)
	if siteURL == "" {
		siteURL = "http://localhost"
	}
	if _, err := licenseSvc.VerifyLicense(siteURL, req.LicenseKey); err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}

	secret := generateSecret()

	// 根据 driver 生成 config 和 DSN
	var configContent string
	var dsn string

	if req.Driver == "sqlite" {
		// SQLite：DBName 是文件路径
		dbPath := req.DBName
		if !filepath.IsAbs(dbPath) {
			// 相对路径：放在当前目录
			dbPath = filepath.Join(".", dbPath)
		}
		// 确保目录存在
		os.MkdirAll(filepath.Dir(dbPath), 0755)
		// SQLite DSN 是完整路径
		dsn = dbPath
		configContent = fmt.Sprintf(`server:
  port: 8080
  mode: release

jwt:
  secret: %s
  expire: 720

database:
  driver: sqlite
  dbname: %s

logger:
  level: info
  path: logs/app.log
`, secret, dbPath)
	} else {
		// MySQL
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			req.DBUser, req.DBPassword, req.DBHost, req.DBPort, req.DBName)
		configContent = fmt.Sprintf(`server:
  port: 8080
  mode: release

jwt:
  secret: %s
  expire: 720

database:
  driver: mysql
  host: "%s"
  port: %d
  user: "%s"
  password: "%s"
  dbname: "%s"

logger:
  level: info
  path: logs/app.log
`, secret, req.DBHost, req.DBPort, req.DBUser, req.DBPassword, req.DBName)
	}

	// 连接数据库
	if err := db.Init(dsn); err != nil {
		if req.Driver == "mysql" {
			// MySQL：尝试创建数据库
			dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
				req.DBUser, req.DBPassword, req.DBHost, req.DBPort)
			if err2 := db.Init(dsnWithoutDB); err2 != nil {
				response.Error(ctx, response.CodeServerError, "数据库连接失败: "+err.Error())
				return
			}
			if err := db.DB.Exec("CREATE DATABASE IF NOT EXISTS `" + req.DBName + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error; err != nil {
				response.Error(ctx, response.CodeServerError, "创建数据库失败: "+err.Error())
				return
			}
			if underlyingDB, err := db.DB.DB(); err == nil && underlyingDB != nil {
				underlyingDB.Close()
			}
			if err := db.Init(dsn); err != nil {
				response.Error(ctx, response.CodeServerError, "连接目标数据库失败: "+err.Error())
				return
			}
		} else {
			response.Error(ctx, response.CodeServerError, "数据库连接失败: "+err.Error())
			return
		}
	}

	// 更新全局 config，以便 db.Init 使用正确驱动
	if config.AppConfig == nil {
		config.AppConfig = &config.Config{}
	}
	if config.AppConfig.Database.Driver == "" {
		config.AppConfig.Database.Driver = req.Driver
	}

	// 执行迁移（完整模型列表，与 cmd/main.go 保持一致）
	if err := db.DB.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.ProductSku{},
		&model.Card{},
		&model.Order{},
		&model.OrderCard{},
		&model.Setting{},
		&model.Coupon{},
		&model.UserCoupon{},
		&model.Seckill{},
		&model.WholesaleRule{},
		&model.Commission{},
		&model.Withdraw{},
		&model.PointsLog{},
		&model.SigninLog{},
		&model.AuditLog{},
		&model.DistributionTree{},
	); err != nil {
		response.Error(ctx, response.CodeServerError, "数据库迁移失败: "+err.Error())
		return
	}

	// 创建或更新管理员
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.AdminPass), bcrypt.DefaultCost)
	if err != nil {
		response.Error(ctx, response.CodeServerError, "密码加密失败: "+err.Error())
		return
	}

	var existingUser model.User
	if err := db.DB.Where("email = ?", req.AdminEmail).First(&existingUser).Error; err == nil {
		// 用户已存在，更新
		db.DB.Model(&model.User{}).Where("email = ?", req.AdminEmail).Updates(map[string]interface{}{
			"password": string(hashed),
			"nickname": req.AdminUser,
			"is_admin": true,
		})
	} else {
		admin := &model.User{
			Email:    req.AdminEmail,
			Password: string(hashed),
			Nickname: req.AdminUser,
			Level:    9,
			IsAdmin:  true,
			Status:   1,
		}
		if err := db.DB.Create(admin).Error; err != nil {
			response.Error(ctx, response.CodeServerError, "创建管理员失败: "+err.Error())
			return
		}
	}

	// 写入配置文件
	if err := os.WriteFile("config.yaml", []byte(configContent), 0644); err != nil {
		response.Error(ctx, response.CodeServerError, "写入配置文件失败: "+err.Error())
		return
	}

	// 写入授权配置文件 config/Auth.php
	if err := licenseSvc.SaveAuthPHP(&service.AuthInfo{
		AppName:     "晨泽发卡",
		AppVersion:  versionString(),
		AuthCode:    req.LicenseKey,
		URL:         siteURL,
		InstalledAt: int64(time.Now().Unix()),
	}); err != nil {
		response.Error(ctx, response.CodeServerError, "写入授权配置失败: "+err.Error())
		return
	}

	// 创建 install.lock
	if err := os.WriteFile("install.lock", []byte("installed"), 0644); err != nil {
		response.Error(ctx, response.CodeServerError, "创建安装锁文件失败: "+err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "安装成功",
		"admin":   req.AdminUser,
	})
}

// generateSecret 生成随机 JWT secret
func generateSecret() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*"
	b := make([]byte, 48)
	for i := range b {
		b[i] = chars[i%len(chars)]
	}
	return string(b)
}

// versionString 返回当前版本号
func versionString() string {
	if appcfg.Version != "" {
		return appcfg.Version
	}
	return "v2.3.0"
}
