package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/e-commerce/shopper/src/config"
)

type Cart struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	UserID    uint       `json:"userId"`
	CartItems []CartItem `gorm:"ForeignKey:CartID" json:"cartItems"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
}

func (cart *Cart) TableName() string {
	return "cart"
}

func CreateCart(userId uint) {
	cart := Cart{UserID: userId, CreatedAt: time.Now()}
	if dbc := DB.Create(&cart); dbc.Error != nil {
		log.Printf("Failed to create cart for user id: %d", userId)
	}
}

func UpdateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Context().Value(config.UserKey).(uint)
	var (
		cart      Cart
		item      CartItem
		cartItems []CartItem
	)
	json.NewDecoder(r.Body).Decode(&item)
	DB.Model(&Cart{}).Where("user_id=?", userId).Find(&cart)
	item.CartID = cart.ID
	if item.Quantity == 0 {
		response := DeleteCartItemHelper(item, userId)
		json.NewEncoder(w).Encode(response)
		return
	}
	if DB.Model(&CartItem{}).Where(
		"product_id=? AND cart_id=?", item.ProductID, item.CartID,
	).Update("quantity", item.Quantity).RowsAffected == 1 {
		msg := fmt.Sprintf("Quantity updated to %d", item.Quantity)
		DB.Model(&cart).Association("CartItems").Find(&cartItems)
		json.NewEncoder(w).Encode(&Response{Message: msg, Data: prepareCartData(cartItems), Error: ""})
		return
	}

	if DB.Create(&item).RowsAffected == 1 {
		DB.Model(&cart).Association("CartItems").Find(&cartItems)
		json.NewEncoder(w).Encode(&Response{Message: "Item added to cart", Data: prepareCartData(cartItems), Error: ""})
	}
}

func DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Context().Value(config.UserKey).(uint)
	var item CartItem
	json.NewDecoder(r.Body).Decode(&item)
	response := DeleteCartItemHelper(item, userId)
	json.NewEncoder(w).Encode(response)
}

func DeleteCartItemHelper(item CartItem, userId uint) (response *Response) {
	var (
		cart      Cart
		cartItems []CartItem
	)
	DB.Model(&Cart{}).Where("user_id=?", userId).Find(&cart)
	dbc := DB.Unscoped().Where("cart_id=? AND product_id=?", item.CartID, item.ProductID).Delete(&CartItem{})
	if dbc.Error != nil {
		response = &Response{Message: "Could not remove item from cart", Data: []string{}, Error: dbc.Error.Error()}
		return
	}
	DB.Model(&cart).Association("CartItems").Find(&cartItems)
	response = &Response{Message: "Item removed from cart.", Data: prepareCartData(cartItems), Error: ""}
	return
}

func GetCartItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Context().Value(config.UserKey).(uint)
	var (
		cart      Cart
		cartItems []CartItem
	)
	DB.Model(&Cart{}).Where("user_id=?", userId).Find(&cart)
	DB.Model(&cart).Association("CartItems").Find(&cartItems)
	json.NewEncoder(w).Encode(&Response{Message: "Cart items fetched successfully", Data: prepareCartData(cartItems), Error: ""})
}

func prepareCartData(cartItems []CartItem) []CartItemRes {
	data := []CartItemRes{}
	for _, item := range cartItems {
		var product Product
		DB.Find(&product, item.ProductID)
		data = append(data, CartItemRes{item, product})
	}
	return data
}
