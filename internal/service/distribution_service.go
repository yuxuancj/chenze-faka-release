package service

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/db"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type DistributionService struct{}

func NewDistributionService() *DistributionService { return &DistributionService{} }

// BindParent 绑定上下级关系（注册时调用）
func (s *DistributionService) BindParent(userID, parentID uint) error {
	if db.DB == nil || parentID == 0 {
		return nil
	}
	if userID == parentID {
		return errors.New("不能绑定自己")
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// 检查是否已有上级
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}
		if user.ParentID > 0 {
			return errors.New("已有上级，不可重复绑定")
		}
		// 检查上级是否是下级（防止循环）
		ancestor := parentID
		for i := 0; i < 5 && ancestor > 0; i++ {
			if ancestor == userID {
				return errors.New("不能绑定自己的下级")
			}
			var p model.User
			if err := tx.Select("parent_id").First(&p, ancestor).Error; err != nil {
				break
			}
			ancestor = p.ParentID
		}
		// 更新 parent_id
		return tx.Model(&user).Update("parent_id", parentID).Error
	})
}

// GetInviteCode 生成/获取用户邀请码
func (s *DistributionService) GetInviteCode(userID uint) (string, error) {
	if db.DB == nil {
		return "", errors.New("数据库未连接")
	}
	var user model.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return "", err
	}
	if user.InviteCode != "" {
		return user.InviteCode, nil
	}
	// 生成6位唯一邀请码
	for i := 0; i < 20; i++ {
		code := randomStr(6)
		var cnt int64
		db.DB.Model(&model.User{}).Where("invite_code=?", code).Count(&cnt)
		if cnt == 0 {
			db.DB.Model(&user).Update("invite_code", code)
			return code, nil
		}
	}
	return fmt.Sprintf("U%d", userID), nil
}

// RecordCommission 订单支付成功后，按三级分销计算佣金
func (s *DistributionService) RecordCommission(tx *gorm.DB, userID uint, orderID uint, amount float64) error {
	if tx == nil {
		return errors.New("事务未初始化")
	}
	ss := NewSettingService()
	if !ss.GetBool("distrib_enabled", true) {
		return nil
	}
	level1Rate := ss.GetFloat("distrib_level1_rate", 10) / 100
	level2Rate := ss.GetFloat("distrib_level2_rate", 5) / 100
	level3Rate := ss.GetFloat("distrib_level3_rate", 2) / 100

	levels := []struct {
		Level   int
		Percent float64
	}{
		{1, level1Rate},
		{2, level2Rate},
		{3, level3Rate},
	}

	currentID := userID
	for _, lv := range levels {
		var user model.User
		if err := tx.Select("parent_id").First(&user, currentID).Error; err != nil {
			return nil
		}
		if user.ParentID == 0 {
			return nil
		}
		commission := amount * lv.Percent
		if commission <= 0 {
			return nil
		}
		rec := &model.Commission{
			UserID:     user.ParentID,
			FromUserID: currentID,
			OrderID:    orderID,
			Amount:     commission,
			Level:      lv.Level,
			Status:     model.CommissionStatusPending,
		}
		if err := tx.Create(rec).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.User{}).Where("id=?", user.ParentID).
			Update("balance", gorm.Expr("balance + ?", commission)).Error; err != nil {
			return err
		}
		currentID = user.ParentID
	}
	return nil
}

// GeneratePoster 生成推广海报
func (s *DistributionService) GeneratePoster(inviteCode string) (string, error) {
	baseURL := "https://api.qrserver.com/v1/create-qr-code/?size=200x200&data="
	inviteLink := "https://example.com/register?invite=" + inviteCode
	qrURL := baseURL + inviteLink
	return qrURL, nil
}

// GetCommissions 查询佣金明细
func (s *DistributionService) GetCommissions(userID uint, page, size int) (int64, []model.Commission, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.Commission
	db.DB.Model(&model.Commission{}).Where("user_id=?", userID).Count(&total)
	if err := db.DB.Offset((page - 1) * size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

// GetSummary 获取分销统计信息
func (s *DistributionService) GetSummary(userID uint) (map[string]interface{}, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	
	var totalCommission float64
	db.DB.Model(&model.Commission{}).Where("user_id=?", userID).Select("IFNULL(SUM(amount),0)").Scan(&totalCommission)
	
	var availableCommission float64
	db.DB.Model(&model.Commission{}).Where("user_id=? AND status=1", userID).Select("IFNULL(SUM(amount),0)").Scan(&availableCommission)
	
	var teamCount int64
	db.DB.Model(&model.User{}).Where("parent_id=?", userID).Count(&teamCount)
	
	code, _ := s.GetInviteCode(userID)
	
	return map[string]interface{}{
		"total_commission":     totalCommission,
		"available_commission": availableCommission,
		"team_count":           teamCount,
		"invite_code":          code,
	}, nil
}

// GetTeam 获取团队成员（直接下级）
func (s *DistributionService) GetTeam(userID uint, page, size int) (int64, []model.User, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.User
	db.DB.Model(&model.User{}).Where("parent_id=?", userID).Count(&total)
	if err := db.DB.Offset((page-1)*size).Limit(size).Order("id desc").Select("id, nickname, created_at, balance, points").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

// ApplyWithdraw 申请提现
func (s *DistributionService) ApplyWithdraw(userID uint, amount float64, accountType, account, realName string) (*model.Withdraw, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	if amount <= 0 {
		return nil, errors.New("提现金额必须大于0")
	}
	w := &model.Withdraw{
		UserID:      userID,
		Amount:      amount,
		AccountType: accountType,
		Account:     account,
		RealName:    realName,
		Status:      model.WithdrawStatusPending,
	}
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}
		if user.Balance < amount {
			return errors.New("余额不足")
		}
		// 冻结余额
		if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}
		return tx.Create(w).Error
	})
	if err != nil {
		return nil, err
	}
	return w, nil
}

// ListWithdraws 提现记录
func (s *DistributionService) ListWithdraws(userID uint, page, size int) (int64, []model.Withdraw, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.Withdraw
	query := db.DB.Model(&model.Withdraw{})
	if userID > 0 {
		query = query.Where("user_id=?", userID)
	}
	query.Count(&total)
	if err := query.Offset((page - 1) * size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

// ProcessWithdraw 后台审核提现
func (s *DistributionService) ProcessWithdraw(id uint, approved bool, remark string) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var w model.Withdraw
		if err := tx.First(&w, id).Error; err != nil {
			return err
		}
		if w.Status != model.WithdrawStatusPending {
			return errors.New("已处理")
		}
		if approved {
			if err := tx.Model(&w).Updates(map[string]interface{}{
				"status": model.WithdrawStatusApproved,
				"remark": remark,
			}).Error; err != nil {
				return err
			}
		} else {
			// 拒绝：退回余额
			if err := tx.Model(&w).Updates(map[string]interface{}{
				"status": model.WithdrawStatusRejected,
				"remark": remark,
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Where("id=?", w.UserID).
				Update("balance", gorm.Expr("balance + ?", w.Amount)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
