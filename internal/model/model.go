package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

const (
	UserStatusActive = 1
	UserStatusBanned = 0

	ProductStatusOn  = 1
	ProductStatusOff = 0

	OrderStatusPending   = 0
	OrderStatusPaid      = 1
	OrderStatusCompleted = 2
	OrderStatusClosed    = 3
	OrderStatusRefunded  = 4

	PayTypeEpay   = "epay"
	PayTypeAlipay = "alipay"
	PayTypeBalance = "balance"

	CouponTypeDiscount = 1
	CouponTypeFull     = 2

	CouponStatusActive = 1
	CouponStatusOff    = 0

	SeckillStatusActive = 1
	SeckillStatusOff    = 0

	CommissionStatusPending  = 0
	CommissionStatusSettled  = 1

	WithdrawStatusPending  = 0
	WithdrawStatusApproved = 1
	WithdrawStatusRejected = 2
)

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSON value")
	}
	return json.Unmarshal(bytes, j)
}

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Email      string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password   string    `gorm:"size:255;not null" json:"-"`
	Nickname   string    `gorm:"size:50" json:"nickname"`
	Avatar     string    `gorm:"size:255" json:"avatar"`
	Balance    float64   `gorm:"type:decimal(10,2);default:0" json:"balance"`
	Points     int       `gorm:"default:0" json:"points"`
	Level      int       `gorm:"default:1" json:"level"`
	IsAdmin    bool      `gorm:"default:false" json:"is_admin"`
	Status     int       `gorm:"default:1" json:"status"`
	ParentID   uint      `gorm:"default:0;index" json:"parent_id"`
	InviteCode string    `gorm:"size:20;uniqueIndex" json:"invite_code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ParentID  uint      `gorm:"default:0" json:"parent_id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `json:"category_id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int       `gorm:"default:0" json:"stock"`
	Sales       int       `gorm:"default:0" json:"sales"`
	Type        string    `gorm:"size:20;default:'card'" json:"type"`
	Image       string    `gorm:"size:255" json:"image"`
	Status      int       `gorm:"default:1" json:"status"`
	IsHidden    bool      `gorm:"default:false" json:"is_hidden"`
	HasSku      bool      `gorm:"default:false" json:"has_sku"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductSku struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProductID  uint      `gorm:"index;not null" json:"product_id"`
	SkuCode    string    `gorm:"size:50" json:"sku_code"`
	SpecNames  JSONMap   `gorm:"type:json" json:"spec_names"`
	Price      float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock      int       `gorm:"default:0" json:"stock"`
	Sales      int       `gorm:"default:0" json:"sales"`
	Image      string    `gorm:"size:255" json:"image"`
	Weight     float64   `gorm:"type:decimal(10,2);default:0" json:"weight"`
	Status     int       `gorm:"default:1" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Card struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	CardData  string    `gorm:"type:text;not null" json:"card_data"`
	Status    int       `gorm:"default:0" json:"status"`
	OrderID   uint      `gorm:"index;default:0" json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	OrderNo         string     `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	UserID          uint       `gorm:"index;not null" json:"user_id"`
	ProductID       uint       `json:"product_id"`
	SkuID           uint       `gorm:"default:0" json:"sku_id"`
	ProductSnapshot string     `gorm:"type:text" json:"product_snapshot"`
	SkuSnapshot     string     `gorm:"type:text" json:"sku_snapshot"`
	Quantity        int        `gorm:"default:1" json:"quantity"`
	Amount          float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	PayType         string     `gorm:"size:20" json:"pay_type"`
	Status          int        `gorm:"default:0;index" json:"status"`
	PaidAt          *time.Time `json:"paid_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	Email           string     `gorm:"size:100" json:"email"`
	Remark          string     `gorm:"size:500" json:"remark"`
	CouponID        uint       `gorm:"default:0" json:"coupon_id"`
	CouponDiscount  float64    `gorm:"type:decimal(10,2);default:0" json:"coupon_discount"`
	PointsUsed      int        `gorm:"default:0" json:"points_used"`
	PointsDiscount  float64    `gorm:"type:decimal(10,2);default:0" json:"points_discount"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type OrderCard struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"index;not null" json:"order_id"`
	CardID    uint      `json:"card_id"`
	CardData  string    `gorm:"type:text" json:"card_data"`
	CreatedAt time.Time `json:"created_at"`
}

type Setting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"size:50;uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Coupon struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Type        int       `gorm:"default:1" json:"type"`
	Value       float64   `gorm:"type:decimal(10,2);not null" json:"value"`
	MinAmount   float64   `gorm:"type:decimal(10,2);default:0" json:"min_amount"`
	Total       int       `gorm:"default:0" json:"total"`
	Used        int       `gorm:"default:0" json:"used"`
	ExpireStart *time.Time `json:"expire_start"`
	ExpireEnd   *time.Time `json:"expire_end"`
	Status      int       `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserCoupon struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	CouponID  uint      `gorm:"index;not null" json:"coupon_id"`
	OrderID   uint      `gorm:"default:0;index" json:"order_id"`
	Status    int       `gorm:"default:0" json:"status"`
	UsedAt    *time.Time `json:"used_at"`
	ExpireAt  *time.Time `json:"expire_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Seckill struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ProductID    uint      `gorm:"index;not null" json:"product_id"`
	SkuID        uint      `gorm:"default:0" json:"sku_id"`
	SeckillPrice float64   `gorm:"type:decimal(10,2);not null" json:"seckill_price"`
	Stock        int       `gorm:"default:0" json:"stock"`
	Sold         int       `gorm:"default:0" json:"sold"`
	PreheatStart *time.Time `json:"preheat_start"`
	PreheatEnd   *time.Time `json:"preheat_end"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	LimitPerUser int       `gorm:"default:1" json:"limit_per_user"`
	Status       int       `gorm:"default:1" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type WholesaleRule struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	ProductID  uint    `gorm:"default:0;index" json:"product_id"`
	MinQty     int     `gorm:"default:1" json:"min_qty"`
	MaxQty     int     `gorm:"default:0" json:"max_qty"`
	DiscountRate float64 `gorm:"type:decimal(5,4);default:1" json:"discount_rate"`
	FixedPrice float64 `gorm:"type:decimal(10,2);default:0" json:"fixed_price"`
	CreatedAt  time.Time `json:"created_at"`
}

type Commission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	FromUserID uint     `gorm:"index;not null" json:"from_user_id"`
	OrderID   uint      `gorm:"index;not null" json:"order_id"`
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Level     int       `gorm:"default:1" json:"level"`
	Status    int       `gorm:"default:0" json:"status"`
	SettledAt *time.Time `json:"settled_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Withdraw struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	AccountType string   `gorm:"size:20" json:"account_type"`
	Account   string    `gorm:"size:100" json:"account"`
	RealName  string    `gorm:"size:50" json:"real_name"`
	Status    int       `gorm:"default:0" json:"status"`
	Remark    string    `gorm:"size:500" json:"remark"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt time.Time `json:"created_at"`
}

type PointsLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Amount      int       `gorm:"not null" json:"amount"`
	Type        string    `gorm:"size:20" json:"type"`
	Description string    `gorm:"size:200" json:"description"`
	OrderID     uint      `gorm:"default:0;index" json:"order_id"`
	BalanceAfter int      `gorm:"default:0" json:"balance_after"`
	CreatedAt   time.Time `json:"created_at"`
}

type SigninLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ContinuousDays int   `gorm:"default:1" json:"continuous_days"`
	RewardPoints int     `gorm:"default:0" json:"reward_points"`
	CreatedAt time.Time `json:"created_at"`
}
