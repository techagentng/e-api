package models

type Order struct {
	Model
	UserID     uint    `json:"user_id"`
	ProductID  uint    `json:"product_id"`
	Quantity   int     `json:"quantity" binding:"required"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status" gorm:"default:'Pending'"` // Statuses: 'Pending', 'Completed', 'Canceled'
}
