package service

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/db"
	"chenze-faka/internal/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// OrderItem 订单项
type OrderItem struct {
	ProductID uint  `json:"product_id"`
	SkuID     uint  `json:"sku_id"`
	Quantity  int   `json:"quantity"`
	Price     float64 `json:"price"`
}

// AdvancedOrderService 增强版订单服务
type AdvancedOrderService struct{}

func NewAdvancedOrderService() *AdvancedOrderService { return &AdvancedOrderService{} }

// CreateAdvanced 创建支持SKU、优惠券、积分、批发折扣的订单
func (s *AdvancedOrderService) CreateAdvanced(userID uint, items []OrderItem, couponID uint, usePoints int, email, payType, remark string) (*model.Order, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	if len(items) == 0 {
		return nil, errors.New("购物车为空")
	}
	var finalOrder *model.Order

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		var primaryItem OrderItem
		// 计算每个商品的价格（应用批发折扣）并锁定库存
		wholesaleSvc := NewWholesaleService()
		skuSvc := NewSkuService()
		itemSnapshots := make([]string, 0, len(items))
		for _, item := range items {
			var product model.Product
			if err := tx.Raw(fmt.Sprintf("SELECT * FROM products WHERE id=%d LIMIT 1 FOR UPDATE", item.ProductID)).First(&product).Error; err != nil {
				return errors.New("商品不存在")
			}
			if product.Status != 1 {
				return errors.New("商品已下架")
			}
			// 应用批发折扣获取最终单价
			unitPrice := product.Price
			if item.SkuID > 0 {
				// SKU商品
				sku, err := skuSvc.ConsumeInTx(tx, item.SkuID, item.Quantity)
				if err != nil {
					return err
				}
				unitPrice = sku.Price
				itemSnapshots = append(itemSnapshots, fmt.Sprintf("SKU-%d:%s x%d", sku.ID, product.Name, item.Quantity))
			} else {
				// 普通商品
				if product.Stock < item.Quantity {
					return errors.New("库存不足: " + product.Name)
				}
				res := tx.Model(&model.Product{}).Where("id=? AND stock >= ?", product.ID, item.Quantity).
					Update("stock", gorm.Expr("stock - ?", item.Quantity))
				if res.Error != nil {
					return res.Error
				}
				if res.RowsAffected == 0 {
					return errors.New("库存不足: " + product.Name)
				}
				itemSnapshots = append(itemSnapshots, fmt.Sprintf("%s x%d", product.Name, item.Quantity))
			}
			// 尝试批发折扣
			finalUnitPrice, _ := wholesaleSvc.ApplyInTx(tx, product.ID, unitPrice, item.Quantity)
			itemTotal := finalUnitPrice * float64(item.Quantity)
			totalAmount += itemTotal

			// 记录第一个商品作为主商品
			if primaryItem.ProductID == 0 {
				primaryItem = item
				primaryItem.Price = finalUnitPrice
			}
		}

		// 应用优惠券
		var couponDiscount float64
		if couponID > 0 {
			couponSvc := NewCouponService()
			discount, err := couponSvc.CalculateDiscount(couponID, totalAmount)
			if err == nil && discount > 0 {
				couponDiscount = discount
			}
		}
		if totalAmount-couponDiscount < 0 {
			couponDiscount = totalAmount
		}

		// 应用积分抵扣（最多抵扣订单金额的50%）
		var pointsDiscount float64
		var pointsUsed int
		if usePoints > 0 && userID > 0 {
			pointsSvc := NewPointsService()
			pointsValue, err := pointsSvc.UsePointsInTx(tx, userID, usePoints, 0, 0.5, totalAmount-couponDiscount)
			if err == nil && pointsValue > 0 {
				pointsDiscount = pointsValue
				pointsUsed = usePoints
			}
		}

		finalAmount := totalAmount - couponDiscount - pointsDiscount
		if finalAmount < 0.01 {
			finalAmount = 0.01
		}

		// 记录使用的优惠券
		if couponID > 0 && couponDiscount > 0 && userID > 0 {
			couponSvc := NewCouponService()
			couponSvc.UseCoupon(tx, userID, couponID, 0) // order_id在创建后更新
		}

		// 创建订单
		orderNo := generateOrderNo()
		snapshotJSON, _ := json.Marshal(items)
		order := &model.Order{
			OrderNo:         orderNo,
			UserID:          userID,
			ProductID:       primaryItem.ProductID,
			SkuID:           primaryItem.SkuID,
			ProductSnapshot: fmt.Sprintf("%s, 明细: %s", totalAmountToDesc(totalAmount, couponDiscount, pointsDiscount), fmtItems(itemSnapshots)),
			SkuSnapshot:     string(snapshotJSON),
			Quantity:        primaryItem.Quantity,
			Amount:          finalAmount,
			PayType:         payType,
			Status:          model.OrderStatusPending,
			Email:           email,
			Remark:          remark,
			CouponID:        couponID,
			CouponDiscount:  couponDiscount,
			PointsUsed:      pointsUsed,
			PointsDiscount:  pointsDiscount,
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		// 更新优惠券的 order_id
		if couponID > 0 && couponDiscount > 0 {
			tx.Model(&model.UserCoupon{}).
				Where("user_id=? AND coupon_id=? AND order_id=0", userID, couponID).
				Update("order_id", order.ID)
		}
		finalOrder = order
		return nil
	})
	if err != nil {
		return nil, err
	}
	return finalOrder, nil
}

// MarkPaidAdvanced 支付成功后的处理（含佣金结算、积分奖励）
func (s *AdvancedOrderService) MarkPaidAdvanced(orderNo, payType string) (*model.Order, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	logger.Infof("支付回调处理: order_no=%s, pay_type=%s", orderNo, payType)
	var order *model.Order
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var lockedOrder model.Order
		rows, err := tx.Raw(fmt.Sprintf("SELECT id FROM orders WHERE order_no='%s' LIMIT 1 FOR UPDATE", orderNo)).Rows()
		if err != nil {
			logger.Errorf("锁定订单失败: %v", err)
			return err
		}
		var id uint
		if rows.Next() {
			rows.Scan(&id)
		}
		rows.Close()
		if id == 0 {
			logger.Warnf("订单不存在: order_no=%s", orderNo)
			return gorm.ErrRecordNotFound
		}
		if err := tx.First(&lockedOrder, id).Error; err != nil {
			return err
		}
		if lockedOrder.Status != model.OrderStatusPending {
			logger.Infof("订单已支付，跳过重复处理: order_no=%s, status=%d", orderNo, lockedOrder.Status)
			order = &lockedOrder
			return nil
		}
		// 3. 卡密发货（若无SKU则按原逻辑）
		if lockedOrder.SkuID == 0 {
			cardSvc := NewCardService()
			_, err := cardSvc.consumeInTx(tx, lockedOrder.ProductID, lockedOrder.Quantity, lockedOrder.ID)
			if err != nil {
				// 卡密发货失败不影响订单状态（可能是无卡密商品）
				_ = err
			}
		}
		// 4. 更新订单状态
		now := lockedOrder.CreatedAt.Add(0) // placeholder
		if err := tx.Model(&lockedOrder).Updates(map[string]interface{}{
			"status":   model.OrderStatusPaid,
			"paid_at":  now,
			"pay_type": payType,
		}).Error; err != nil {
			return err
		}
		// 5. 增加商品销量
		if err := tx.Model(&model.Product{}).Where("id=?", lockedOrder.ProductID).
			Update("sales", gorm.Expr("sales + ?", lockedOrder.Quantity)).Error; err != nil {
			return err
		}
		// 6. 完成订单
		if err := tx.Model(&lockedOrder).Updates(map[string]interface{}{
			"status":       model.OrderStatusCompleted,
			"completed_at": now,
		}).Error; err != nil {
			return err
		}
		// 7. 积分奖励（消费1元=1积分）
		if lockedOrder.UserID > 0 {
			pointsGain := int(lockedOrder.Amount)
			if pointsGain > 0 {
				pointsSvc := NewPointsService()
				pointsSvc.AddPoints(lockedOrder.UserID, pointsGain, "purchase",
					fmt.Sprintf("订单消费奖励 %.2f元", lockedOrder.Amount), lockedOrder.ID)
			}
			// 8. 分销佣金结算
			distSvc := NewDistributionService()
			distSvc.RecordCommission(tx, lockedOrder.UserID, lockedOrder.ID, lockedOrder.Amount)
		}
		order = &lockedOrder
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

// 辅助：格式化商品明细
func fmtItems(items []string) string {
	if len(items) == 0 {
		return ""
	}
	result := ""
	for i, it := range items {
		if i > 0 {
			result += "; "
		}
		result += it
	}
	return result
}

func totalAmountToDesc(total, coupon, points float64) string {
	desc := fmt.Sprintf("总额%.2f", total)
	if coupon > 0 {
		desc += fmt.Sprintf(" 优惠券-%.2f", coupon)
	}
	if points > 0 {
		desc += fmt.Sprintf(" 积分-%.2f", points)
	}
	return desc
}
