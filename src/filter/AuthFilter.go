package filter

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/e-commerce/shopper/src/auth"
	"github.com/e-commerce/shopper/src/config"
)

var PUBLIC_PATHS = []string{"/login", "/user/register", "/products", "/viber", "/send"}

func Auth() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			for _, p := range PUBLIC_PATHS {
				if strings.Contains(path, p) {
					h.ServeHTTP(w, r)
					return
				}
			}
			value := r.Header.Get("Authorization")
			if value == "" || !strings.HasPrefix(value, "Bearer ") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			value = strings.Split(value, "Bearer ")[1]

			// Initialize a new instance of `Claims`
			claims := &auth.Claims{}

			// Parse the JWT string and store the result in `claims`.
			// Note that we are passing the key in this method as well. This method will return an error
			// if the token is invalid (if it has expired according to the expiry time we set on sign in),
			// or if the signature does not match
			token, err := jwt.ParseWithClaims(value, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte("SECRET_KEY"), nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), config.UserKey, claims.UserId)
			req := r.WithContext(ctx)
			h.ServeHTTP(w, req)
		})
	}
}
