package model

import "time"

type Payment struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"-"`
	Method     string    `json:"method"`
	UPI        string    `json:"upi"`
	CardOwner  string    `json:"cardOwner"`
	CardNumber uint64    `json:"cardNumber"`
	ExpiryDate string    `json:"expiryDate"`
	UserID     uint      `json:"-"`
	User       User      `json:"-"`
}

func (payment *Payment) TableName() string {
	return "payment"
}
