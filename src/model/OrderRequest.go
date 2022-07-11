package model

type OrderRequest struct {
	OrderItems []OrderItem `json:"items"`
	Payment    Payment     `json:"payment"`
	Address    Address     `json:"address"`
}
