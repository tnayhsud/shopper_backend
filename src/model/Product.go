package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Product struct {
	ID            uint    `gorm:"primarykey" json:"id"`
	ImgUrl        string  `gorm:"type:VARCHAR(255);NOT NULL" json:"imgUrl"`
	Discount      float64 `gorm:"type:INT;default:0" json:"discount"`
	Gender        string  `gorm:"type:VARCHAR(16);NOT NULL" json:"gender"`
	Title         string  `gorm:"type:VARCHAR(255);NOT NULL" json:"title"`
	DiscountPrice float64 `gorm:"type:DECIMAL(7,2);NOT NULL" json:"discountPrice"`
	ActualPrice   float64 `gorm:"type:DECIMAL(7,2);NOT NULL" json:"actualPrice"`
}

func (product *Product) TableName() string {
	return "product"
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": "Products saved successfully",
	}
	var products []Product
	json.NewDecoder(r.Body).Decode(&products)
	if dbc := DB.Create(&products); dbc.Error != nil {
		log.Printf("Error: %s\n", dbc.Error)
		response = map[string]string{
			"message": "Failed to save products",
		}
	} else {
		log.Println("Products saved successfully")
	}
	json.NewEncoder(w).Encode(response)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var product Product
	if dbc := DB.First(&product, params["id"]); dbc.Error != nil {
		log.Printf("Error: %s\n", dbc.Error)
		response := map[string]string{
			"message": fmt.Sprintf("No product found for id: %s", params["id"]),
		}
		json.NewEncoder(w).Encode(response)
	} else {
		json.NewEncoder(w).Encode(&product)
	}
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var products []Product
	gender := r.URL.Query().Get("gender")
	if gender == "" || gender == "All" {
		DB.Order("rand()").Find(&products)
	} else {
		DB.Where("gender = ?", gender).Find(&products)
	}
	json.NewEncoder(w).Encode(&products)
}
