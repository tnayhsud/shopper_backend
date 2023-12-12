package model

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/e-commerce/shopper/src/config"
)

type Order struct {
	ID         uint        `gorm:"primarykey" json:"id"`
	UserID     uint        `json:"userId"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"orderItems"`
	PaymentID  uint        `json:"paymentId"`
	AddressID  uint        `json:"addressId"`
	CreatedAt  time.Time   `gorm:"autoCreateTime" json:"createdAt"`
}

func (order *Order) TableName() string {
	return "order"
}

func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Context().Value(config.UserKey).(uint)
	var orderRequest OrderRequest
	json.NewDecoder(r.Body).Decode(&orderRequest)
	order := &Order{UserID: userId}
	if dbc := DB.Create(&order); dbc.Error != nil {
		handleError(w, "Unknown error occured", dbc.Error.Error())
		return
	}
	for _, item := range orderRequest.OrderItems {
		item.OrderID = order.ID
		DB.Create(&item)
	}
	address := orderRequest.Address
	if dbc := DB.Find(&Address{}, address.ID); dbc.RowsAffected == 0 {
		address.UserID = userId
		if dbc := DB.Create(&address); dbc.Error != nil {
			handleError(w, "Unknown error occured", dbc.Error.Error())
			return
		}
	}
	if address.IsDefault {
		var addressList []Address
		DB.Model(&Address{}).Where("user_id", userId).Find(&addressList)
		for i := range addressList {
			if addressList[i].IsDefault && addressList[i].ID != address.ID {
				DB.Model(&addressList[i]).Update("is_default", false)
			}
		}
	}
	payment := orderRequest.Payment
	payment.UserID = userId
	if dbc := DB.Create(&payment); dbc.Error != nil {
		handleError(w, "Payment could not be completed", dbc.Error.Error())
		return
	}
	order.AddressID = address.ID
	order.PaymentID = payment.ID
	DB.Model(&Order{}).Where("id=?", order.ID).Updates(&order)

	for _, item := range orderRequest.OrderItems {
		DB.Where("product_id=?", item.ProductID).Delete(&CartItem{})
	}
	json.NewEncoder(w).Encode(&Response{"Order placed successfully", map[string]uint{"orderId": order.ID}, ""})
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Context().Value(config.UserKey).(uint)
	var orders []Order
	DB.Model(&Order{}).Where("user_id=?", userId).Find(&orders)
	for i := range orders {
		var orderItems []OrderItem
		DB.Model(&orders[i]).Association("OrderItems").Find(&orderItems)
		orders[i].OrderItems = prepareOrderItemData(orderItems)
	}
	json.NewEncoder(w).Encode(&Response{"Orders fetched successfully", orders, ""})
}

func prepareOrderItemData(items []OrderItem) []OrderItem {
	for i := range items {
		var product Product
		DB.Find(&product, items[i].ProductID)
		items[i].Product = product
	}
	return items
}

func handleError(w http.ResponseWriter, msg, err string) {
	log.Printf("%s : %s", msg, err)
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(&Response{"Order could not be placed", "", msg})
}
