package service

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/db"
	"chenze-faka/internal/pkg/logger"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SeckillService struct{}

func NewSeckillService() *SeckillService { return &SeckillService{} }

// List 后台列表
func (s *SeckillService) List(page, size int) (int64, []model.Seckill, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.Seckill
	db.DB.Model(&model.Seckill{}).Count(&total)
	if err := db.DB.Offset((page - 1) * size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

// ActiveList 前台活动列表
func (s *SeckillService) ActiveList() ([]model.Seckill, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var list []model.Seckill
	now := time.Now()
	// 预热中 or 进行中
	if err := db.DB.Where("status=1 AND stock > 0 AND (preheat_start <= ? OR start_time <= ?)", now, now).
		Where("(end_time IS NULL OR end_time >= ?)", now).
		Order("start_time asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetByID 获取单个活动
func (s *SeckillService) GetByID(id uint) (*model.Seckill, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	sk := &model.Seckill{}
	if err := db.DB.First(sk, id).Error; err != nil {
		return nil, err
	}
	return sk, nil
}

// Create 创建活动
func (s *SeckillService) Create(sk *model.Seckill) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	if sk.Status == 0 {
		sk.Status = 1
	}
	return db.DB.Create(sk).Error
}

// Update 更新
func (s *SeckillService) Update(id uint, data map[string]interface{}) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Model(&model.Seckill{}).Where("id=?", id).Updates(data).Error
}

// Delete 删除
func (s *SeckillService) Delete(id uint) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Delete(&model.Seckill{}, id).Error
}

// Order 秒杀下单（行锁保证并发安全）
func (s *SeckillService) Order(seckillID, userID uint, qty int, email, payType, remark string) (*model.Order, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	logger.Infof("秒杀下单请求: seckill_id=%d, user_id=%d, qty=%d", seckillID, userID, qty)
	var order *model.Order
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var sk model.Seckill
		if err := tx.Raw(fmt.Sprintf("SELECT * FROM seckills WHERE id=%d LIMIT 1 FOR UPDATE", seckillID)).First(&sk).Error; err != nil {
			logger.Warnf("秒杀活动不存在: seckill_id=%d", seckillID)
			return errors.New("秒杀活动不存在")
		}
		if sk.Status != 1 {
			return errors.New("活动已结束")
		}
		now := time.Now()
		// 2. 检查时间
		if sk.StartTime != nil && now.Before(*sk.StartTime) {
			return errors.New("活动尚未开始")
		}
		if sk.EndTime != nil && now.After(*sk.EndTime) {
			return errors.New("活动已结束")
		}
		// 3. 检查库存
		if sk.Stock < qty {
			return errors.New("库存不足")
		}
		// 4. 检查每人限购
		var bought int64
		tx.Model(&model.Order{}).Where("product_id=? AND user_id=? AND status >= ?", sk.ProductID, userID, model.OrderStatusPending).
			Select("IFNULL(SUM(quantity),0)").Scan(&bought)
		if int64(sk.LimitPerUser) > 0 && int64(sk.LimitPerUser) < bought+int64(qty) {
			return errors.New("超出限购数量")
		}
		// 5. 锁定商品
		var product model.Product
		if err := tx.First(&product, sk.ProductID).Error; err != nil {
			return errors.New("商品不存在")
		}
		// 6. 扣减秒杀库存（原子递减）
		res := tx.Model(&model.Seckill{}).Where("id=? AND stock >= ?", sk.ID, qty).
			Update("stock", gorm.Expr("stock - ?", qty))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("库存不足")
		}
		// 更新已售
		if err := tx.Model(&model.Seckill{}).Where("id=?", sk.ID).
			Update("sold", gorm.Expr("sold + ?", qty)).Error; err != nil {
			return err
		}
		// 7. 创建订单
		orderNo := generateOrderNo()
		amount := sk.SeckillPrice * float64(qty)
		skuID := sk.SkuID
		snapshot := ""
		if sk.SkuID > 0 {
			var sku model.ProductSku
			if err := tx.First(&sku, sk.SkuID).Error; err == nil {
				snapshot = fmt.Sprintf("%s|规格:ID-%d|秒杀价:%.2f", product.Name, sku.ID, sk.SeckillPrice)
			}
		} else {
			snapshot = fmt.Sprintf("%s|秒杀价:%.2f", product.Name, sk.SeckillPrice)
		}
		order = &model.Order{
			OrderNo:         orderNo,
			UserID:          userID,
			ProductID:       sk.ProductID,
			SkuID:           skuID,
			ProductSnapshot: snapshot,
			SkuSnapshot:     snapshot,
			Quantity:        qty,
			Amount:          amount,
			PayType:         payType,
			Status:          model.OrderStatusPending,
			Email:           email,
			Remark:          remark,
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}
