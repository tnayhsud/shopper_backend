package model

import (
	"time"

	"gorm.io/gorm"
)

type CartItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CartID    uint           `json:"cartId"`
	Cart      Cart           `json:"-"`
	ProductID uint           `json:"productId"`
	Quantity  uint           `gorm:"type:INT;NOT NULL" json:"quantity"`
}

type CartItemRes struct {
	CartItem
	Product Product `json:"product"`
}

func (cartItem *CartItem) TableName() string {
	return "cart_item"
}
