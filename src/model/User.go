package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/e-commerce/shopper/src/config"
	"github.com/gorilla/mux"
)

type User struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"type:VARCHAR(128);NOT NULL" json:"name"`
	Email       string    `gorm:"type:VARCHAR(255);NOT NULL" json:"email"`
	Password    string    `gorm:"type:VARCHAR(64);NOT NULL" json:"password"`
	Payment     []Payment `gorm:"ForeignKey:UserID" json:"paymentList"`
	AddressList []Address `gorm:"ForeignKey:UserID" json:"addressList"`
}

func (user *User) TableName() string {
	return "user"
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": "User saved successfully",
	}
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	if dbc := DB.Create(&user); dbc.Error != nil {
		log.Printf("Error: %s\n", dbc.Error)
		response = map[string]string{
			"message": "Failed to save user",
		}
	} else {
		log.Println("User saved successfully")
		CreateCart(user.ID)
		CreateWishlist(user.ID)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	if dbc := DB.First(&user, params["id"]); dbc.Error != nil {
		log.Printf("Error: %s\n", dbc.Error)
		response := map[string]string{
			"message": fmt.Sprintf("No user found for id: %s", params["id"]),
		}
		json.NewEncoder(w).Encode(response)
	} else {
		json.NewEncoder(w).Encode(&user)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	DB.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

func GetAddressList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Context().Value(config.UserKey).(uint)
	var addressList []Address
	DB.Model(&Address{}).Where("user_id=?", userId).Find(&addressList)
	json.NewEncoder(w).Encode(&addressList)
}
