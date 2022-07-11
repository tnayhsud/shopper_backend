package model

import (
	"time"

	"gorm.io/gorm"
)

type WishlistItem struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	WishlistID uint           `json:"wishlistId"`
	Wishlist   Wishlist       `json:"-"`
	ProductID  uint           `json:"productId"`
	Quantity   uint           `gorm:"type:INT;NOT NULL" json:"quantity"`
}

type WishlistItemRes struct {
	WishlistItem
	Product Product `json:"product"`
}

func (wishlistItem *WishlistItem) TableName() string {
	return "wishlist_item"
}
