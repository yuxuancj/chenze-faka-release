package service

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/db"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CouponService struct{}

func NewCouponService() *CouponService { return &CouponService{} }

// List 后台列表
func (s *CouponService) List(page, size int) (int64, []model.Coupon, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.Coupon
	db.DB.Model(&model.Coupon{}).Count(&total)
	if err := db.DB.Offset((page - 1) * size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

// Available 前台可领取列表
func (s *CouponService) Available() ([]model.Coupon, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	now := time.Now()
	var list []model.Coupon
	query := db.DB.Model(&model.Coupon{}).Where("status=1 AND used < total")
	query = query.Where("(expire_start IS NULL OR expire_start <= ?) AND (expire_end IS NULL OR expire_end >= ?)", now, now)
	if err := query.Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// Create 创建优惠券
func (s *CouponService) Create(c *model.Coupon) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	if c.Status == 0 {
		c.Status = 1
	}
	return db.DB.Create(c).Error
}

// Update 更新
func (s *CouponService) Update(id uint, data map[string]interface{}) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Model(&model.Coupon{}).Where("id=?", id).Updates(data).Error
}

// Delete 删除
func (s *CouponService) Delete(id uint) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Delete(&model.Coupon{}, id).Error
}

// Claim 兑换/领取优惠券
func (s *CouponService) Claim(userID, couponID uint) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var coupon model.Coupon
		if err := tx.First(&coupon, couponID).Error; err != nil {
			return errors.New("优惠券不存在")
		}
		if coupon.Status != 1 {
			return errors.New("优惠券不可用")
		}
		if coupon.Total > 0 && coupon.Used >= coupon.Total {
			return errors.New("优惠券已领完")
		}
		now := time.Now()
		if coupon.ExpireStart != nil && now.Before(*coupon.ExpireStart) {
			return errors.New("优惠券未到领取时间")
		}
		if coupon.ExpireEnd != nil && now.After(*coupon.ExpireEnd) {
			return errors.New("优惠券已过期")
		}
		// 每人限领1张检查
		var userCount int64
		tx.Model(&model.UserCoupon{}).Where("user_id=? AND coupon_id=?", userID, couponID).Count(&userCount)
		if userCount > 0 {
			return errors.New("已领取过该优惠券")
		}
		// 创建用户优惠券
		uc := &model.UserCoupon{
			UserID:   userID,
			CouponID: couponID,
			Status:   0,
			ExpireAt: coupon.ExpireEnd,
		}
		if err := tx.Create(uc).Error; err != nil {
			return err
		}
		// 更新优惠券已用数
		return tx.Model(&coupon).Update("used", gorm.Expr("used + 1")).Error
	})
}

// UserCoupons 我的优惠券
func (s *CouponService) UserCoupons(userID uint, status int) ([]model.UserCoupon, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var list []model.UserCoupon
	query := db.DB.Where("user_id=?", userID)
	if status >= 0 {
		query = query.Where("status=?", status)
	}
	if err := query.Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// CalculateDiscount 计算优惠金额
func (s *CouponService) CalculateDiscount(couponID uint, amount float64) (float64, error) {
	if db.DB == nil {
		return 0, errors.New("数据库未连接")
	}
	var coupon model.Coupon
	if err := db.DB.First(&coupon, couponID).Error; err != nil {
		return 0, errors.New("优惠券不存在")
	}
	if coupon.Status != 1 {
		return 0, errors.New("优惠券不可用")
	}
	if amount < coupon.MinAmount {
		return 0, errors.New("订单金额未达到优惠券最低使用金额")
	}
	switch coupon.Type {
	case 1: // 满减
		return coupon.Value, nil
	case 2: // 折扣 (value表示折扣率, 如 0.8 = 8折)
		return amount * (1 - coupon.Value), nil
	default:
		return 0, errors.New("无效的优惠券类型")
	}
}

// UseCoupon 使用优惠券
func (s *CouponService) UseCoupon(tx *gorm.DB, userID, couponID, orderID uint) (float64, error) {
	var uc model.UserCoupon
	if err := tx.Where("user_id=? AND coupon_id=? AND status=0", userID, couponID).
		First(&uc).Error; err != nil {
		return 0, errors.New("优惠券不存在或已使用")
	}
	now := time.Now()
	if uc.ExpireAt != nil && now.After(*uc.ExpireAt) {
		return 0, errors.New("优惠券已过期")
	}
	// 标记为已使用
	if err := tx.Model(&uc).Updates(map[string]interface{}{
		"status":   1,
		"order_id": orderID,
		"used_at":  &now,
	}).Error; err != nil {
		return 0, err
	}
	return 0, nil // 金额由Controller计算
}
