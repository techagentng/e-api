package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	ProductID  uint      `json:"product_id" gorm:"not null"`
	Quantity   int       `json:"quantity" binding:"required"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status" gorm:"default:'Pending'"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Items      []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	User    User    `json:"user" gorm:"foreignKey:UserID"`
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

type OrderItem struct {
    ID         uint    `json:"id" gorm:"primaryKey"`
    OrderID    uint `json:"order_id" gorm:"type:uuid"`
    ProductID  uint    `json:"product_id" binding:"required"`
    Quantity   int     `json:"quantity" binding:"required"`
    UnitPrice  float64 `json:"unit_price"`
    TotalPrice float64 `json:"total_price"`
}

type OrderRequest struct {
    Items []OrderItem `json:"items"` 
}