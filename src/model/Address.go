package model

type Address struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	AddressLine string `json:"addressLine"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	ZipCode     uint   `json:"zipCode"`
	IsDefault   bool   `json:"isDefault"`
	UserID      uint   `json:"userId"`
	User        User   `json:"-"`
}

func (address *Address) TableName() string {
	return "address"
}
