package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/e-commerce/shopper/src/config"
	"github.com/e-commerce/shopper/src/model"
)

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Claims struct {
	UserId uint   `json:"userId"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbc := model.DB.Model(&model.User{}).Where("email=?", creds.Email).Find(&user)
	if dbc.Error != nil || user.Password != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := generateJWT(user, []byte("SECRET_KEY"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", token)
	json.NewEncoder(w).Encode(&user)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Context().Value(config.UserKey).(uint)
	var user model.User
	model.DB.Model(&model.User{}).Find(&user, userId)
	token, err := generateJWT(user, []byte("SECRET_KEY"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", token)
	json.NewEncoder(w).Encode(&user)
}

func generateJWT(user model.User, secret []byte) (signedToken string, err error) {
	log.Println("Generating JWT")
	claims := &Claims{
		UserId: user.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString(secret)
	return
}

func expirationTime() time.Time {
	return time.Now().Add(120 * time.Minute)
}
