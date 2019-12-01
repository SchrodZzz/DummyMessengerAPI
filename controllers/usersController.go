package controllers

import (
	"DummyMessengerAPI/models"
	u "DummyMessengerAPI/utils"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Response(w, http.StatusForbidden, u.Message(false, "Invalid request"))
		return
	}
	u.Response(w, http.StatusOK, user.Create())
}

var Authorize = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Response(w, http.StatusForbidden, u.Message(false, "Invalid request"))
		return
	}
	u.Response(w, http.StatusOK, models.Login(user.Login, user.Password))
}

var Logout = func(w http.ResponseWriter, r *http.Request) {
	token, _ := jwt.Parse(strings.Split(r.Header.Get("Authorization"), " ")[1],
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
	u.Response(w, http.StatusOK, models.Logout(token.Raw))
}
