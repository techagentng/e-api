package models

type Product struct {
	ID        uint      `gorm:"primaryKey"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Quantity    int     `json:"quantity" binding:"required"`
	Orders   []Order   `json:"orders" gorm:"foreignKey:ProductID"`
}
