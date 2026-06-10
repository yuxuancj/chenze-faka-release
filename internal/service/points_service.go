package service

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/db"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PointsService struct{}

func NewPointsService() *PointsService { return &PointsService{} }

// AddPoints 增加积分（消费、签到等）
func (s *PointsService) AddPoints(userID uint, amount int, typ, desc string, orderID uint) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	if amount <= 0 {
		return nil
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}
		newBalance := user.Points + amount
		if err := tx.Model(&user).Update("points", newBalance).Error; err != nil {
			return err
		}
		return tx.Create(&model.PointsLog{
			UserID:      userID,
			Amount:      amount,
			Type:        typ,
			Description: desc,
			OrderID:     orderID,
			BalanceAfter: newBalance,
		}).Error
	})
}

// UsePoints 使用积分抵扣
func (s *PointsService) UsePoints(userID uint, points int, orderID uint, maxDiscountRatio float64, orderAmount float64) (float64, error) {
	if db.DB == nil {
		return 0, errors.New("数据库未连接")
	}
	ss := NewSettingService()
	if !ss.GetBool("points_enabled", true) {
		return 0, nil
	}
	if points <= 0 {
		return 0, nil
	}
	var discount float64
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}
		if user.Points < points {
			return errors.New("积分不足")
		}
		pointsRate := ss.GetFloat("points_discount_rate", 1)
		pointsValue := float64(points) * pointsRate / 100.0
		effectiveRatio := maxDiscountRatio
		if effectiveRatio <= 0 {
			effectiveRatio = ss.GetFloat("points_max_discount_percent", 50) / 100
		}
		maxDiscount := orderAmount * effectiveRatio
		if pointsValue > maxDiscount {
			pointsValue = maxDiscount
		}
		discount = pointsValue
		newBalance := user.Points - points
		if err := tx.Model(&user).Update("points", newBalance).Error; err != nil {
			return err
		}
		return tx.Create(&model.PointsLog{
			UserID:       userID,
			Amount:       -points,
			Type:         "consume",
			Description:  fmt.Sprintf("订单消费抵扣 %.2f元", pointsValue),
			OrderID:      orderID,
			BalanceAfter: newBalance,
		}).Error
	})
	return discount, err
}

// UsePointsInTx 在事务内使用积分
func (s *PointsService) UsePointsInTx(tx *gorm.DB, userID uint, points int, orderID uint, maxDiscountRatio float64, orderAmount float64) (float64, error) {
	if points <= 0 {
		return 0, nil
	}
	ss := NewSettingService()
	if !ss.GetBool("points_enabled", true) {
		return 0, nil
	}
	var user model.User
	if err := tx.First(&user, userID).Error; err != nil {
		return 0, err
	}
	if user.Points < points {
		return 0, errors.New("积分不足")
	}
	pointsRate := ss.GetFloat("points_discount_rate", 1)
	pointsValue := float64(points) * pointsRate / 100.0
	effectiveRatio := maxDiscountRatio
	if effectiveRatio <= 0 {
		effectiveRatio = ss.GetFloat("points_max_discount_percent", 50) / 100
	}
	maxDiscount := orderAmount * effectiveRatio
	if pointsValue > maxDiscount {
		pointsValue = maxDiscount
	}
	newBalance := user.Points - points
	if err := tx.Model(&user).Update("points", newBalance).Error; err != nil {
		return 0, err
	}
	if err := tx.Create(&model.PointsLog{
		UserID:       userID,
		Amount:       -points,
		Type:         "consume",
		Description:  fmt.Sprintf("订单消费抵扣 %.2f元", pointsValue),
		OrderID:      orderID,
		BalanceAfter: newBalance,
	}).Error; err != nil {
		return 0, err
	}
	return pointsValue, nil
}

// GetBalance 获取积分余额
func (s *PointsService) GetBalance(userID uint) (int, error) {
	if db.DB == nil {
		return 0, errors.New("数据库未连接")
	}
	var user model.User
	if err := db.DB.Select("points").First(&user, userID).Error; err != nil {
		return 0, err
	}
	return user.Points, nil
}

// GetLogs 积分明细
func (s *PointsService) GetLogs(userID uint, page, size int) (int64, []model.PointsLog, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.PointsLog
	query := db.DB.Model(&model.PointsLog{})
	if userID > 0 {
		query = query.Where("user_id=?", userID)
	}
	query.Count(&total)
	if err := query.Offset((page - 1) * size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

// SignIn 每日签到
func (s *PointsService) SignIn(userID uint) (int, int, error) {
	if db.DB == nil {
		return 0, 0, errors.New("数据库未连接")
	}
	ss := NewSettingService()
	if !ss.GetBool("points_enabled", true) {
		return 0, 0, errors.New("积分功能已关闭")
	}
	var rewardPoints, continuousDays int
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		var todayCount int64
		tx.Model(&model.SigninLog{}).Where("user_id=? AND created_at >= ?", userID, todayStart).Count(&todayCount)
		if todayCount > 0 {
			return errors.New("今日已签到")
		}
		yesterdayStart := todayStart.AddDate(0, 0, -1)
		var lastSign model.SigninLog
		err := tx.Where("user_id=? AND created_at >= ? AND created_at < ?", userID, yesterdayStart, todayStart).
			Order("created_at desc").First(&lastSign).Error
		if err == nil {
			continuousDays = lastSign.ContinuousDays + 1
		} else {
			continuousDays = 1
		}
		rewardPoints = ss.GetInt("points_signin_reward", 10)
		if continuousDays >= 7 {
			rewardPoints = rewardPoints * 2
		} else if continuousDays >= 3 {
			rewardPoints += ss.GetInt("points_continuous_reward", 5)
		}
		// 创建签到记录
		if err := tx.Create(&model.SigninLog{
			UserID:         userID,
			ContinuousDays: continuousDays,
			RewardPoints:   rewardPoints,
		}).Error; err != nil {
			return err
		}
		// 更新积分
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}
		newBalance := user.Points + rewardPoints
		if err := tx.Model(&user).Update("points", newBalance).Error; err != nil {
			return err
		}
		// 写入积分日志
		return tx.Create(&model.PointsLog{
			UserID:       userID,
			Amount:       rewardPoints,
			Type:         "signin",
			Description:  fmt.Sprintf("连续签到%d天", continuousDays),
			BalanceAfter: newBalance,
		}).Error
	})
	return rewardPoints, continuousDays, err
}

// CanSignIn 检查今日是否可签到
func (s *PointsService) CanSignIn(userID uint) bool {
	if db.DB == nil {
		return false
	}
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var count int64
	db.DB.Model(&model.SigninLog{}).Where("user_id=? AND created_at >= ?", userID, todayStart).Count(&count)
	return count == 0
}

// GetContinuousDays 获取当前连续签到天数
func (s *PointsService) GetContinuousDays(userID uint) int {
	if db.DB == nil {
		return 0
	}
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	var log model.SigninLog
	err := db.DB.Where("user_id=? AND created_at >= ?", userID, yesterdayStart).
		Order("created_at desc").First(&log).Error
	if err != nil {
		return 0
	}
	// 今日已签到则用当前记录，否则昨日的记录+1
	if log.CreatedAt.After(todayStart) {
		return log.ContinuousDays
	}
	return log.ContinuousDays + 1
}
