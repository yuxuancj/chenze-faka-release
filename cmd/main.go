package main

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/config"
	"chenze-faka/internal/pkg/db"
	"chenze-faka/internal/pkg/logger"
	"chenze-faka/internal/router"
	"chenze-faka/internal/service"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := logger.Init("storage/logs"); err != nil {
		fmt.Printf("logger init failed: %v\n", err)
	}
	logger.Infof("starting chenze-faka...")

	// 配置加载
	cfgPath := "config.yaml"
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		if _, err := os.Stat("config.yaml.example"); err == nil {
			cfgPath = "config.yaml.example"
		} else {
			logger.Fatalf("config.yaml not found")
		}
	}
	if err := config.Load(cfgPath); err != nil {
		logger.Fatalf("load config failed: %v", err)
	}

	// 数据库初始化
	if err := db.Init(config.AppConfig.Database.DSN()); err != nil {
		logger.Errorf("db init failed: %v (will continue with degraded mode)", err)
	} else {
		if err := db.DB.AutoMigrate(
			&model.User{},
			&model.Category{},
			&model.Product{},
			&model.Card{},
			&model.Order{},
			&model.OrderCard{},
			&model.Setting{},
		); err != nil {
			logger.Errorf("auto migrate failed: %v", err)
		} else {
			logger.Infof("db migrated")
			seedIfEmpty()
		}
	}

	os.MkdirAll("storage/backups", 0755)
	os.MkdirAll("storage/uploads", 0755)

	if config.AppConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 加载模板：先扫描 html 文件，再逐个注册为其相对路径
	baseDir := "templates"
	if info, err := os.Stat(baseDir); err == nil && info.IsDir() {
		htmlFiles := collectHTML(baseDir)
		if len(htmlFiles) > 0 {
			templ := gin.New().HTMLRender
			_ = templ
			// 直接用 LoadHTMLGlob 做初步加载，再覆盖
			r.LoadHTMLGlob(filepath.Join(baseDir, "**/*"))
		}
	}

	router.Setup(r)

	port := config.AppConfig.Server.Port
	if port == 0 {
		port = 8080
	}
	logger.Infof("server listening on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}

func collectHTML(baseDir string) []string {
	var files []string
	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func seedIfEmpty() {
	if db.DB == nil {
		return
	}
	var userCount int64
	db.DB.Model(&model.User{}).Count(&userCount)
	if userCount == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := &model.User{
			Email:    "admin@chenze.com",
			Password: string(hashed),
			Nickname: "超级管理员",
			Level:    9,
			IsAdmin:  true,
			Status:   1,
		}
		if err := db.DB.Create(admin).Error; err != nil {
			logger.Errorf("seed admin failed: %v", err)
		} else {
			logger.Infof("default admin created: admin@chenze.com / admin123")
		}
	}

	var productCount int64
	db.DB.Model(&model.Product{}).Count(&productCount)
	if productCount == 0 {
		cat := &model.Category{Name: "默认分类", Sort: 1, Status: 1}
		db.DB.Create(cat)
		p := &model.Product{
			CategoryID:  cat.ID,
			Name:        "示例商品 - VIP 会员",
			Description: "这是一个示例商品，用于演示发卡系统。",
			Price:       9.9,
			Stock:       10,
			Type:        "card",
			Status:      1,
		}
		db.DB.Create(p)
		cards := []string{"示例卡密1 - CODE001", "示例卡密2 - CODE002", "示例卡密3 - CODE003"}
		_, _ = service.NewCardService().Import(p.ID, cards)
		logger.Infof("sample product created: id=%d", p.ID)
	}
}
