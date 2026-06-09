package model

import "time"

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

	PayTypeEpay  = "epay"
	PayTypeAlipay = "alipay"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Nickname  string    `gorm:"size:50" json:"nickname"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Balance   float64   `gorm:"type:decimal(10,2);default:0" json:"balance"`
	Points    int       `gorm:"default:0" json:"points"`
	Level     int       `gorm:"default:1" json:"level"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Type        string    `gorm:"size:20;default:'card'" json:"type"` // card / text / api
	Image       string    `gorm:"size:255" json:"image"`
	Status      int       `gorm:"default:1" json:"status"`
	IsHidden    bool      `gorm:"default:false" json:"is_hidden"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Card struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	CardData  string    `gorm:"type:text;not null" json:"card_data"`
	Status    int       `gorm:"default:0" json:"status"` // 0=未使用 1=已售出
	OrderID   uint      `gorm:"index;default:0" json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderNo    string    `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	ProductID  uint      `gorm:"not null" json:"product_id"`
	ProductSnapshot string `gorm:"type:text" json:"product_snapshot"`
	Quantity   int       `gorm:"default:1" json:"quantity"`
	Amount     float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	PayType    string    `gorm:"size:20" json:"pay_type"`
	Status     int       `gorm:"default:0;index" json:"status"`
	PaidAt     *time.Time `json:"paid_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Email      string    `gorm:"size:100" json:"email"`
	Remark     string    `gorm:"size:500" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
