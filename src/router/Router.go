package router

import (
	"github.com/e-commerce/shopper/src/auth"
	model "github.com/e-commerce/shopper/src/model"
	"github.com/gorilla/mux"
)

var Router mux.Router

func init() {
	Router = *mux.NewRouter()

	// Authentication routes
	Router.HandleFunc("/login", auth.Authenticate).Methods("POST")
	Router.HandleFunc("/refresh/token", auth.RefreshToken).Methods("GET")

	// Product routes
	Router.HandleFunc("/products", model.CreateProduct).Methods("POST")
	Router.HandleFunc("/products", model.GetProducts).Methods("GET")
	Router.HandleFunc("/products/{id}", model.GetProduct).Methods("GET")

	// User routes
	Router.HandleFunc("/user/register", model.CreateUser).Methods("POST")
	Router.HandleFunc("/users", model.GetUsers).Methods("GET")
	// Router.HandleFunc("/user/{id}", model.GetUser).Methods("GET")
	Router.HandleFunc("/user/address", model.GetAddressList).Methods("GET")

	// Cart routes
	Router.HandleFunc("/cart/items", model.UpdateCart).Methods("POST")
	Router.HandleFunc("/cart/items", model.GetCartItems).Methods("GET")

	// Wishlist routes
	Router.HandleFunc("/wishlist/items", model.UpdateWishlist).Methods("POST")
	Router.HandleFunc("/wishlist/items", model.GetWishlistItems).Methods("GET")

	// Order routes
	Router.HandleFunc("/order", model.PlaceOrder).Methods("POST")
	Router.HandleFunc("/orders", model.GetOrders).Methods("GET")

	Router.HandleFunc("/viber", model.Viber).Methods("POST")
	Router.HandleFunc("/send", model.Send).Methods("POST")
}
