package model

import (
	"fmt"

	"github.com/e-commerce/shopper/src/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	user := config.ReadConfig("datasource.user")
	pwd := config.ReadConfig("datasource.pwd")
	host := config.ReadConfig("datasource.host")
	schema := config.ReadConfig("datasource.schema")
	queryParams := "charset=utf8mb4&parseTime=True&loc=Local"
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", user, pwd, host, schema, queryParams)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		panic("Failure to connect to database")
	}
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Cart{})
	db.AutoMigrate(&CartItem{})
	db.AutoMigrate(&Wishlist{})
	db.AutoMigrate(&WishlistItem{})
	db.AutoMigrate(&Address{})
	db.AutoMigrate(&Payment{})
	db.AutoMigrate(&OrderItem{})
	db.AutoMigrate(&Order{})
	DB = db
}
