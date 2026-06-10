package service

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/db"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService { return &UserService{} }

func (s *UserService) Register(email, password, nickname string) (*model.User, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	email = strings.ToLower(strings.TrimSpace(email))
	var count int64
	db.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return nil, errors.New("邮箱已注册")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if nickname == "" {
		nickname = "用户" + randomStr(6)
	}
	u := &model.User{
		Email:    email,
		Password: string(hashed),
		Nickname: nickname,
		Level:    1,
		Status:   model.UserStatusActive,
	}
	if err := db.DB.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) Login(email, password string) (*model.User, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	email = strings.ToLower(strings.TrimSpace(email))
	u := &model.User{}
	if err := db.DB.Where("email = ?", email).First(u).Error; err != nil {
		return nil, errors.New("邮箱或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("邮箱或密码错误")
	}
	if u.Status != model.UserStatusActive {
		return nil, errors.New("账号已被禁用")
	}
	return u, nil
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	u := &model.User{}
	if err := db.DB.First(u, id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) UpdateProfile(id uint, nickname string) error {
	if !db.IsReady() {
		return errors.New("数据库未连接")
	}
	return db.DB.Model(&model.User{}).Where("id=?", id).Update("nickname", nickname).Error
}

// UpdateUserByAdmin 管理员更新用户（余额、积分、等级、状态）
func (s *UserService) UpdateUserByAdmin(id uint, updates map[string]interface{}) error {
	if !db.IsReady() {
		return errors.New("数据库未连接")
	}
	return db.DB.Model(&model.User{}).Where("id=?", id).Updates(updates).Error
}

func (s *UserService) ChangePassword(id uint, oldPwd, newPwd string) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	u := &model.User{}
	if err := db.DB.First(u, id).Error; err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPwd)); err != nil {
		return errors.New("原密码错误")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return db.DB.Model(u).Update("password", string(hashed)).Error
}

func (s *UserService) List(page, size int) (int64, []model.User, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.User
	db.DB.Model(&model.User{}).Count(&total)
	if err := db.DB.Offset((page - 1) * size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

type CategoryService struct{}

func NewCategoryService() *CategoryService { return &CategoryService{} }

func (s *CategoryService) All() ([]model.Category, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var list []model.Category
	if err := db.DB.Order("sort asc, id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *CategoryService) Create(name string, parentID uint, sort int) (*model.Category, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	c := &model.Category{Name: name, ParentID: parentID, Sort: sort, Status: 1}
	if err := db.DB.Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) Update(id uint, name string, sort int) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Model(&model.Category{}).Where("id=?", id).Updates(map[string]interface{}{"name": name, "sort": sort}).Error
}

func (s *CategoryService) Delete(id uint) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Delete(&model.Category{}, id).Error
}

type ProductService struct{}

func NewProductService() *ProductService { return &ProductService{} }

func (s *ProductService) List(page, size int, categoryID uint, keyword string) (int64, []model.Product, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.Product
	query := db.DB.Model(&model.Product{}).Where("status=? AND is_hidden=?", 1, false)
	if categoryID > 0 {
		query = query.Where("category_id=?", categoryID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	query.Count(&total)
	if err := query.Offset((page - 1) * size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

func (s *ProductService) AdminList(page, size int, keyword string) (int64, []model.Product, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.Product
	query := db.DB.Model(&model.Product{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	query.Count(&total)
	if err := query.Offset((page - 1) * size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

func (s *ProductService) GetByID(id uint) (*model.Product, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	p := &model.Product{}
	if err := db.DB.First(p, id).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) Create(p *model.Product) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Create(p).Error
}

func (s *ProductService) Update(id uint, data map[string]interface{}) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Model(&model.Product{}).Where("id=?", id).Updates(data).Error
}

func (s *ProductService) Delete(id uint) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Product{}, id).Error; err != nil {
			return err
		}
		return tx.Where("product_id=?", id).Delete(&model.Card{}).Error
	})
}

type CardService struct{}

func NewCardService() *CardService { return &CardService{} }

func (s *CardService) Import(productID uint, cards []string) (int, error) {
	if db.DB == nil {
		return 0, errors.New("数据库未连接")
	}
	tx := db.DB.Begin()
	count := 0
	for _, cd := range cards {
		cd = strings.TrimSpace(cd)
		if cd == "" {
			continue
		}
		if err := tx.Create(&model.Card{ProductID: productID, CardData: cd}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		count++
	}
	if err := tx.Model(&model.Product{}).Where("id=?", productID).
		Update("stock", gorm.Expr("stock + ?", count)).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return count, nil
}

func (s *CardService) ListByProduct(productID uint, page, size int) (int64, []model.Card, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.Card
	db.DB.Model(&model.Card{}).Where("product_id=?", productID).Count(&total)
	if err := db.DB.Offset((page-1)*size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

func (s *CardService) Consume(productID uint, qty int, orderID uint) ([]model.Card, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var cards []model.Card
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 使用 FOR UPDATE（MySQL）或 SQLite 事务隔离保证行级锁，避免并发超卖同一张卡
		rows, err := tx.Raw(consumeLockSQL(productID, qty)).Rows()
		if err != nil {
			return err
		}
		defer rows.Close()
		ids := make([]uint, 0, qty)
		for rows.Next() {
			var id uint
			if err := rows.Scan(&id); err != nil {
				return err
			}
			ids = append(ids, id)
		}
		if len(ids) < qty {
			return errors.New("库存不足")
		}
		// 再次锁定后更新（保证无其他事务抢占）
		if err := tx.Model(&model.Card{}).Where("id IN ?", ids).
			Updates(map[string]interface{}{"status": 1, "order_id": orderID}).Error; err != nil {
			return err
		}
		// 原子递减库存（使用条件更新避免负数）
		res := tx.Model(&model.Product{}).Where("id=? AND stock >= ?", productID, qty).
			Update("stock", gorm.Expr("stock - ?", qty))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("库存不足")
		}
		// 重新加载卡密信息
		if err := tx.Where("id IN ?", ids).Find(&cards).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cards, nil
}

// consumeLockSQL 生成带行锁的 SQL（MySQL FOR UPDATE）
func consumeLockSQL(productID uint, qty int) string {
	return fmt.Sprintf("SELECT id FROM cards WHERE product_id=%d AND status=0 ORDER BY id ASC LIMIT %d FOR UPDATE", productID, qty)
}

type OrderService struct{}

func NewOrderService() *OrderService { return &OrderService{} }

func (s *OrderService) Create(userID, productID uint, qty int, email, payType, remark string) (*model.Order, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var order *model.Order
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 对商品行加排他锁（FOR UPDATE），防止并发下单时超卖
		var product model.Product
		if err := tx.Raw(buildProductLockSQL(productID)).First(&product).Error; err != nil {
			return errors.New("商品不存在")
		}
		if product.Status != 1 {
			return errors.New("商品已下架")
		}
		if product.Stock < qty {
			return errors.New("库存不足")
		}
		// 2. 原子递减商品库存（条件更新：stock >= qty，避免负数）
		res := tx.Model(&model.Product{}).Where("id=? AND stock >= ?", productID, qty).
			Update("stock", gorm.Expr("stock - ?", qty))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("库存不足")
		}
		// 3. 创建订单（状态为 pending，支付后由回调更新）
		orderNo := generateOrderNo()
		amount := product.Price * float64(qty)
		order = &model.Order{
			OrderNo:         orderNo,
			UserID:          userID,
			ProductID:       productID,
			ProductSnapshot: fmt.Sprintf("%s|单价:%.2f", product.Name, product.Price),
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

// buildProductLockSQL 构造商品行锁 SQL（MySQL FOR UPDATE）
func buildProductLockSQL(productID uint) string {
	return fmt.Sprintf("SELECT * FROM products WHERE id=%d LIMIT 1 FOR UPDATE", productID)
}

func (s *OrderService) GetByOrderNo(orderNo string) (*model.Order, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	o := &model.Order{}
	if err := db.DB.Where("order_no=?", orderNo).First(o).Error; err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) GetByID(id uint) (*model.Order, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	o := &model.Order{}
	if err := db.DB.First(o, id).Error; err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) ListByUser(userID uint, page, size int) (int64, []model.Order, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.Order
	db.DB.Model(&model.Order{}).Where("user_id=?", userID).Count(&total)
	if err := db.DB.Offset((page-1)*size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

func (s *OrderService) AdminList(page, size int, keyword string) (int64, []model.Order, error) {
	if db.DB == nil {
		return 0, nil, errors.New("数据库未连接")
	}
	var total int64
	var list []model.Order
	query := db.DB.Model(&model.Order{})
	if keyword != "" {
		query = query.Where("order_no LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)
	if err := query.Offset((page-1)*size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

func (s *OrderService) MarkPaid(orderNo, payType string) (*model.Order, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var order model.Order

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 行锁锁定订单（MySQL: SELECT ... FOR UPDATE，SQLite 事务本身即持有锁）
		var lockedOrder model.Order
		sql := buildLockOrderSQL(orderNo)
		rows, err := tx.Raw(sql).Rows()
		if err != nil {
			return err
		}
		if !rows.Next() {
			rows.Close()
			return gorm.ErrRecordNotFound
		}
		var id uint
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		rows.Close()
		if err := tx.First(&lockedOrder, id).Error; err != nil {
			return err
		}
		order = lockedOrder
		// 2. 幂等性：若已支付/已完成，直接返回成功（避免重复发货）
		if order.Status != model.OrderStatusPending {
			return nil
		}
		// 3. 扣减卡密 + 更新库存（内部已含行锁/事务）
		cards, err := NewCardService().consumeInTx(tx, order.ProductID, order.Quantity, order.ID)
		if err != nil {
			return err
		}
		now := time.Now()
		// 4. 更新订单状态为已支付
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status":   model.OrderStatusPaid,
			"paid_at":  &now,
			"pay_type": payType,
		}).Error; err != nil {
			return err
		}
		// 5. 写入 order_cards 发货记录
		for _, c := range cards {
			if err := tx.Create(&model.OrderCard{OrderID: order.ID, CardID: c.ID, CardData: c.CardData}).Error; err != nil {
				return err
			}
		}
		// 6. 增加商品销量
		if err := tx.Model(&model.Product{}).Where("id=?", order.ProductID).
			Update("sales", gorm.Expr("sales + ?", order.Quantity)).Error; err != nil {
			return err
		}
		// 7. 标记订单已完成
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status":       model.OrderStatusCompleted,
			"completed_at": &now,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// buildLockOrderSQL 构造订单行锁查询 SQL（MySQL FOR UPDATE）
func buildLockOrderSQL(orderNo string) string {
	return fmt.Sprintf("SELECT id FROM orders WHERE order_no='%s' LIMIT 1 FOR UPDATE", orderNo)
}

// consumeInTx 在现有事务内执行卡密扣减（用于 MarkPaid 的链式事务）
func (s *CardService) consumeInTx(tx *gorm.DB, productID uint, qty int, orderID uint) ([]model.Card, error) {
	rows, err := tx.Raw(consumeLockSQL(productID, qty)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := make([]uint, 0, qty)
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if len(ids) < qty {
		return nil, errors.New("库存不足")
	}
	if err := tx.Model(&model.Card{}).Where("id IN ?", ids).
		Updates(map[string]interface{}{"status": 1, "order_id": orderID}).Error; err != nil {
		return nil, err
	}
	res := tx.Model(&model.Product{}).Where("id=? AND stock >= ?", productID, qty).
		Update("stock", gorm.Expr("stock - ?", qty))
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("库存不足")
	}
	var cards []model.Card
	if err := tx.Where("id IN ?", ids).Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *OrderService) GetOrderCards(orderID uint) ([]model.OrderCard, error) {
	if db.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var list []model.OrderCard
	if err := db.DB.Where("order_id=?", orderID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

type SettingService struct{}

func NewSettingService() *SettingService { return &SettingService{} }

// Get 安全读取配置项，db 未初始化或不可用时返回默认值
func (s *SettingService) Get(key, def string) string {
	if !db.IsReady() {
		return def
	}
	setting := &model.Setting{}
	if err := db.DB.Where("`key`=?", key).First(setting).Error; err != nil {
		return def
	}
	return setting.Value
}

func (s *SettingService) Set(key, value string) error {
	if db.DB == nil {
		return errors.New("数据库未连接")
	}
	setting := &model.Setting{}
	err := db.DB.Where("`key`=?", key).First(setting).Error
	if err == gorm.ErrRecordNotFound {
		return db.DB.Create(&model.Setting{Key: key, Value: value}).Error
	}
	if err != nil {
		return err
	}
	return db.DB.Model(setting).Update("value", value).Error
}

func randomStr(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func generateOrderNo() string {
	now := time.Now()
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%d%d", now.UnixNano(), rand.Int63())))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil))[:16])
}
