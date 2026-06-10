package service

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/db"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type SkuService struct{}

func NewSkuService() *SkuService { return &SkuService{} }

// GetSkusByProduct 获取商品的所有 SKU
func (s *SkuService) GetSkusByProduct(productID uint) ([]model.ProductSku, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var list []model.ProductSku
	if err := db.DB.Where("product_id=?", productID).Order("id asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetByID 根据ID获取SKU
func (s *SkuService) GetByID(id uint) (*model.ProductSku, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	sku := &model.ProductSku{}
	if err := db.DB.First(sku, id).Error; err != nil {
		return nil, err
	}
	return sku, nil
}

// ReplaceSkus 替换商品的所有 SKU（先删后增）
func (s *SkuService) ReplaceSkus(productID uint, skus []model.ProductSku) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除所有旧 SKU
		if err := tx.Where("product_id=?", productID).Delete(&model.ProductSku{}).Error; err != nil {
			return err
		}
		// 2. 如果有新SKU，批量插入并更新商品 has_sku 标记
		if len(skus) > 0 {
			for i := range skus {
				skus[i].ProductID = productID
			}
			if err := tx.Create(&skus).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.Product{}).Where("id=?", productID).Update("has_sku", true).Error; err != nil {
				return err
			}
		} else {
			// 无SKU，关闭商品的 SKU 模式
			if err := tx.Model(&model.Product{}).Where("id=?", productID).Update("has_sku", false).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DecrementStock 扣减 SKU 库存（行锁+条件更新，保证并发安全）
func (s *SkuService) DecrementStock(skuID uint, qty int) (*model.ProductSku, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var sku model.ProductSku
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 行锁
		if err := tx.Raw(fmt.Sprintf("SELECT * FROM product_skus WHERE id=%d LIMIT 1 FOR UPDATE", skuID)).First(&sku).Error; err != nil {
			return errors.New("SKU不存在")
		}
		if sku.Status != 1 {
			return errors.New("SKU已下架")
		}
		if sku.Stock < qty {
			return errors.New("SKU库存不足")
		}
		// 原子递减
		res := tx.Model(&model.ProductSku{}).Where("id=? AND stock >= ?", skuID, qty).
			Update("stock", gorm.Expr("stock - ?", qty))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("SKU库存不足")
		}
		// 更新销量
		if err := tx.Model(&model.ProductSku{}).Where("id=?", skuID).
			Update("sales", gorm.Expr("sales + ?", qty)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

// ConsumeInTx 在事务内扣减 SKU 库存（用于下单链式事务）
func (s *SkuService) ConsumeInTx(tx *gorm.DB, skuID uint, qty int) (*model.ProductSku, error) {
	var sku model.ProductSku
	if err := tx.Raw(fmt.Sprintf("SELECT * FROM product_skus WHERE id=%d LIMIT 1 FOR UPDATE", skuID)).First(&sku).Error; err != nil {
		return nil, errors.New("SKU不存在")
	}
	if sku.Status != 1 {
		return nil, errors.New("SKU已下架")
	}
	if sku.Stock < qty {
		return nil, errors.New("SKU库存不足")
	}
	res := tx.Model(&model.ProductSku{}).Where("id=? AND stock >= ?", skuID, qty).
		Update("stock", gorm.Expr("stock - ?", qty))
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("SKU库存不足")
	}
	if err := tx.Model(&model.ProductSku{}).Where("id=?", skuID).
		Update("sales", gorm.Expr("sales + ?", qty)).Error; err != nil {
		return nil, err
	}
	return &sku, nil
}
