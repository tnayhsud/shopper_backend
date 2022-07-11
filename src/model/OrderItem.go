package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	OrderID   uint           `json:"orderId"`
	Order     Order          `json:"-"`
	ProductID uint           `json:"productId"`
	Product   Product        `gorm:"-" json:"product"`
	Quantity  uint           `gorm:"type:INT;NOT NULL" json:"quantity"`
}

func (orderItem *OrderItem) TableName() string {
	return "order_item"
}
