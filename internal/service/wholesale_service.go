package service

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/db"
	"errors"

	"gorm.io/gorm"
)

type WholesaleService struct{}

func NewWholesaleService() *WholesaleService { return &WholesaleService{} }

// List 获取批发规则
func (s *WholesaleService) List(productID uint) ([]model.WholesaleRule, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var list []model.WholesaleRule
	query := db.DB.Model(&model.WholesaleRule{})
	if productID > 0 {
		query = query.Where("product_id=? OR product_id=0", productID)
	}
	if err := query.Order("min_qty asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// Create 创建规则
func (s *WholesaleService) Create(r *model.WholesaleRule) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Create(r).Error
}

// Update 更新
func (s *WholesaleService) Update(id uint, data map[string]interface{}) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Model(&model.WholesaleRule{}).Where("id=?", id).Updates(data).Error
}

// Delete 删除
func (s *WholesaleService) Delete(id uint) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Delete(&model.WholesaleRule{}, id).Error
}

// MatchAndApply 匹配并计算优惠价（返回最终单价和折扣规则）
func (s *WholesaleService) MatchAndApply(productID uint, unitPrice float64, qty int) (finalPrice float64, rule *model.WholesaleRule, err error) {
	if db.DB == nil {
		return unitPrice, nil, nil
	}
	// 先查商品专属规则
	var rules []model.WholesaleRule
	if err := db.DB.Where("product_id=?", productID).Order("min_qty desc").Find(&rules).Error; err != nil {
		return unitPrice, nil, err
	}
	if len(rules) == 0 {
		// 查全局规则
		if err := db.DB.Where("product_id=0").Order("min_qty desc").Find(&rules).Error; err != nil {
			return unitPrice, nil, err
		}
	}
	// 匹配最高满足条件的规则
	for _, r := range rules {
		if qty >= r.MinQty && (r.MaxQty == 0 || qty <= r.MaxQty) {
			if r.FixedPrice > 0 {
				return r.FixedPrice, &r, nil
			}
			if r.DiscountRate > 0 && r.DiscountRate < 1 {
				return unitPrice * r.DiscountRate, &r, nil
			}
		}
	}
	return unitPrice, nil, nil
}

// ApplyInTx 在事务内应用批发折扣（用于下单时）
func (s *WholesaleService) ApplyInTx(tx *gorm.DB, productID uint, unitPrice float64, qty int) (float64, error) {
	var rules []model.WholesaleRule
	if err := tx.Where("product_id=?", productID).Order("min_qty desc").Find(&rules).Error; err != nil {
		return unitPrice, err
	}
	if len(rules) == 0 {
		if err := tx.Where("product_id=0").Order("min_qty desc").Find(&rules).Error; err != nil {
			return unitPrice, err
		}
	}
	for _, r := range rules {
		if qty >= r.MinQty && (r.MaxQty == 0 || qty <= r.MaxQty) {
			if r.FixedPrice > 0 {
				return r.FixedPrice, nil
			}
			if r.DiscountRate > 0 && r.DiscountRate < 1 {
				return unitPrice * r.DiscountRate, nil
			}
		}
	}
	return unitPrice, nil
}
