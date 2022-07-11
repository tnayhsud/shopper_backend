package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/e-commerce/shopper/src/config"
	"github.com/e-commerce/shopper/src/filter"
	_ "github.com/e-commerce/shopper/src/model"
	"github.com/e-commerce/shopper/src/router"
)

func Adapt(h http.Handler, adapters ...filter.Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func main() {
	port := os.Args[1]
	log.Printf("Starting server at %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), Adapt(&router.Router, filter.Auth(), filter.Cors())))
}
