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
	u := &model.User{}
	if err := db.DB.First(u, id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) UpdateProfile(id uint, nickname string) error {
	return db.DB.Model(&model.User{}).Where("id=?", id).Update("nickname", nickname).Error
}

func (s *UserService) ChangePassword(id uint, oldPwd, newPwd string) error {
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
	var list []model.Category
	if err := db.DB.Order("sort asc, id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *CategoryService) Create(name string, parentID uint, sort int) (*model.Category, error) {
	c := &model.Category{Name: name, ParentID: parentID, Sort: sort, Status: 1}
	if err := db.DB.Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) Update(id uint, name string, sort int) error {
	return db.DB.Model(&model.Category{}).Where("id=?", id).Updates(map[string]interface{}{"name": name, "sort": sort}).Error
}

func (s *CategoryService) Delete(id uint) error {
	return db.DB.Delete(&model.Category{}, id).Error
}

type ProductService struct{}

func NewProductService() *ProductService { return &ProductService{} }

func (s *ProductService) List(page, size int, categoryID uint, keyword string) (int64, []model.Product, error) {
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
	p := &model.Product{}
	if err := db.DB.First(p, id).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) Create(p *model.Product) error {
	return db.DB.Create(p).Error
}

func (s *ProductService) Update(id uint, data map[string]interface{}) error {
	return db.DB.Model(&model.Product{}).Where("id=?", id).Updates(data).Error
}

func (s *ProductService) Delete(id uint) error {
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
	var total int64
	var list []model.Card
	db.DB.Model(&model.Card{}).Where("product_id=?", productID).Count(&total)
	if err := db.DB.Offset((page-1)*size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

func (s *CardService) Consume(productID uint, qty int, orderID uint) ([]model.Card, error) {
	var cards []model.Card
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id=? AND status=?", productID, 0).
			Limit(qty).Order("id asc").Find(&cards).Error; err != nil {
			return err
		}
		if len(cards) < qty {
			return errors.New("库存不足")
		}
		ids := make([]uint, 0, len(cards))
		for _, c := range cards {
			ids = append(ids, c.ID)
		}
		if err := tx.Model(&model.Card{}).Where("id IN ?", ids).
			Updates(map[string]interface{}{"status": 1, "order_id": orderID}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Product{}).Where("id=?", productID).
			Update("stock", gorm.Expr("stock - ?", qty)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cards, nil
}

type OrderService struct{}

func NewOrderService() *OrderService { return &OrderService{} }

func (s *OrderService) Create(userID, productID uint, qty int, email, payType, remark string) (*model.Order, error) {
	product, err := NewProductService().GetByID(productID)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	if product.Status != 1 {
		return nil, errors.New("商品已下架")
	}
	if product.Stock < qty {
		return nil, errors.New("库存不足")
	}
	amount := product.Price * float64(qty)
	orderNo := generateOrderNo()
	order := &model.Order{
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
	if err := db.DB.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetByOrderNo(orderNo string) (*model.Order, error) {
	o := &model.Order{}
	if err := db.DB.Where("order_no=?", orderNo).First(o).Error; err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) GetByID(id uint) (*model.Order, error) {
	o := &model.Order{}
	if err := db.DB.First(o, id).Error; err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) ListByUser(userID uint, page, size int) (int64, []model.Order, error) {
	var total int64
	var list []model.Order
	db.DB.Model(&model.Order{}).Where("user_id=?", userID).Count(&total)
	if err := db.DB.Offset((page-1)*size).Limit(size).Order("id desc").Find(&list).Error; err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

func (s *OrderService) AdminList(page, size int, keyword string) (int64, []model.Order, error) {
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
	order, err := s.GetByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	if order.Status != model.OrderStatusPending {
		return order, nil
	}
	cards, err := NewCardService().Consume(order.ProductID, order.Quantity, order.ID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	tx := db.DB.Begin()
	if err := tx.Model(order).Updates(map[string]interface{}{
		"status":  model.OrderStatusPaid,
		"paid_at": &now,
		"pay_type": payType,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, c := range cards {
		if err := tx.Create(&model.OrderCard{OrderID: order.ID, CardID: c.ID, CardData: c.CardData}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Model(&model.Product{}).Where("id=?", order.ProductID).
		Update("sales", gorm.Expr("sales + ?", order.Quantity)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	completed := time.Now()
	if err := tx.Model(order).Updates(map[string]interface{}{
		"status":       model.OrderStatusCompleted,
		"completed_at": &completed,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return order, nil
}

func (s *OrderService) GetOrderCards(orderID uint) ([]model.OrderCard, error) {
	var list []model.OrderCard
	if err := db.DB.Where("order_id=?", orderID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

type SettingService struct{}

func NewSettingService() *SettingService { return &SettingService{} }

func (s *SettingService) Get(key, def string) string {
	setting := &model.Setting{}
	if err := db.DB.Where("`key`=?", key).First(setting).Error; err != nil {
		return def
	}
	return setting.Value
}

func (s *SettingService) Set(key, value string) error {
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
