package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/e-commerce/shopper/src/config"
)

type Wishlist struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UserID        uint           `json:"userId"`
	WishlistItems []WishlistItem `gorm:"ForeignKey:WishlistID" json:"wishlistItems"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"createdAt"`
}

func (wishlist *Wishlist) TableName() string {
	return "wishlist"
}

func CreateWishlist(userId uint) {
	wishlist := Wishlist{UserID: userId, CreatedAt: time.Now()}
	if dbc := DB.Create(&wishlist); dbc.Error != nil {
		log.Printf("Failed to create wishlist for user id: %d", userId)
	}
}

func UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Context().Value(config.UserKey).(uint)
	var (
		wishlist      Wishlist
		item          WishlistItem
		wishlistItems []WishlistItem
	)
	json.NewDecoder(r.Body).Decode(&item)
	DB.Model(&Wishlist{}).Where("user_id=?", userId).Find(&wishlist)
	item.WishlistID = wishlist.ID
	if item.Quantity == 0 {
		response := deleteWishlistItemHelper(item, userId)
		json.NewEncoder(w).Encode(response)
		return
	}
	if DB.Model(&WishlistItem{}).Where(
		"product_id=? AND wishlist_id=?", item.ProductID, item.WishlistID,
	).Update("quantity", item.Quantity).RowsAffected == 1 {
		msg := fmt.Sprintf("Quantity updated to %d", item.Quantity)
		DB.Model(&wishlist).Association("WishlistItems").Find(&wishlistItems)
		json.NewEncoder(w).Encode(&Response{Message: msg, Data: prepareWishlistData(wishlistItems), Error: ""})
		return
	}
	if DB.Create(&item).RowsAffected == 1 {
		DB.Model(&wishlist).Association("WishlistItems").Find(&wishlistItems)
		json.NewEncoder(w).Encode(&Response{Message: "Item added to wishlist", Data: prepareWishlistData(wishlistItems), Error: ""})
	}
}

func DeleteWishlistItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Context().Value(config.UserKey).(uint)
	var item WishlistItem
	json.NewDecoder(r.Body).Decode(&item)
	response := deleteWishlistItemHelper(item, userId)
	json.NewEncoder(w).Encode(response)
}

func deleteWishlistItemHelper(item WishlistItem, userId uint) (response *Response) {
	var (
		wishlist      Wishlist
		wishlistItems []WishlistItem
	)
	DB.Model(&Wishlist{}).Where("user_id=?", userId).Find(&wishlist)
	dbc := DB.Unscoped().Where("wishlist_id=? AND product_id=?", item.WishlistID, item.ProductID).Delete(&WishlistItem{})
	if dbc.Error != nil {
		response = &Response{Message: "Could not remove item from wishlist", Data: []string{}, Error: dbc.Error.Error()}
		return
	}
	DB.Model(&wishlist).Association("WishlistItems").Find(&wishlistItems)
	response = &Response{Message: "Item removed from wishlist.", Data: prepareWishlistData(wishlistItems), Error: ""}
	return
}

func GetWishlistItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Context().Value(config.UserKey).(uint)
	var (
		wishlist      Wishlist
		wishlistItems []WishlistItem
	)
	DB.Model(&Wishlist{}).Where("user_id=?", userId).Find(&wishlist)
	DB.Model(&wishlist).Association("WishlistItems").Find(&wishlistItems)
	json.NewEncoder(w).Encode(&Response{Message: "Wishlist items fetched successfully", Data: prepareWishlistData(wishlistItems), Error: ""})
}

func prepareWishlistData(wsItems []WishlistItem) []WishlistItemRes {
	data := []WishlistItemRes{}
	for _, item := range wsItems {
		var product Product
		DB.Find(&product, item.ProductID)
		data = append(data, WishlistItemRes{item, product})
	}
	return data
}
