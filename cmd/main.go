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
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := logger.Init("storage/logs"); err != nil {
		fmt.Printf("logger init failed: %v\n", err)
	}
	logger.Infof("starting chenze-faka...")

	// 检查是否已安装（install.lock 存在则跳过配置加载和数据库初始化）
	installed := false
	if _, err := os.Stat("install.lock"); err == nil {
		installed = true
	}

	if installed {
		// 已安装：正常加载配置和初始化数据库
		cfgPath := "config.yaml"
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			logger.Fatalf("config.yaml not found")
		}
		if err := config.Load(cfgPath); err != nil {
			logger.Fatalf("load config failed: %v", err)
		}
		if err := db.Init(config.AppConfig.Database.DSN()); err != nil {
			logger.Errorf("db init failed: %v", err)
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
	} else {
		logger.Infof("system not installed, running in install mode")
	}

	os.MkdirAll("storage/backups", 0755)
	os.MkdirAll("storage/uploads", 0755)

	if config.AppConfig != nil && config.AppConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 路由（含 embed 模板加载，运行时无需外部 templates/ 目录）
	router.Setup(r)

	if installed {
		// 已安装：启动超时订单关闭定时任务
		go startOrderExpirer()
	}

	port := 8080
	if config.AppConfig != nil && config.AppConfig.Server.Port > 0 {
		port = config.AppConfig.Server.Port
	}
	logger.Infof("server listening on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}

// startOrderExpirer 启动超时订单关闭定时任务
// 每分钟扫描一次：创建超过 30 分钟且状态仍为 pending 的订单标记为 closed，并恢复商品库存
func startOrderExpirer() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		closeExpiredOrders()
	}
}

// closeExpiredOrders 关闭超时未支付的订单，释放已扣减的库存
func closeExpiredOrders() {
	if !db.IsReady() {
		return
	}
	cutoff := time.Now().Add(-30 * time.Minute)
	var expiredOrders []model.Order
	err := db.DB.Where("status = ? AND created_at < ?", model.OrderStatusPending, cutoff).Find(&expiredOrders).Error
	if err != nil || len(expiredOrders) == 0 {
		return
	}
	for _, order := range expiredOrders {
		err := db.DB.Transaction(func(tx *gorm.DB) error {
			// 再次确认订单状态（防止并发修改）
			var o model.Order
			if err := tx.First(&o, order.ID).Error; err != nil {
				return err
			}
			if o.Status != model.OrderStatusPending {
				return nil
			}
			// 标记订单为已关闭
			if err := tx.Model(&o).Update("status", model.OrderStatusClosed).Error; err != nil {
				return err
			}
			// 恢复商品库存（使用原子条件更新避免重复扣减）
			res := tx.Model(&model.Product{}).Where("id=?", o.ProductID).
				Update("stock", gorm.Expr("stock + ?", o.Quantity))
			if res.Error != nil {
				return res.Error
			}
			logger.Infof("order expired, closed: order_no=%s, restored %d stock", o.OrderNo, o.Quantity)
			return nil
		})
		if err != nil {
			logger.Errorf("close expired order failed: %v", err)
		}
	}
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
