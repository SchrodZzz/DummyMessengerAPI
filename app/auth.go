package app

import (
	"DummyMessengerAPI/models"
	u "DummyMessengerAPI/utils"
	"context"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
	"time"
)

var JwtAuthorization = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/user/new", "/api/user/login"}
		requestPath := r.URL.Path
		for _, path := range notAuth {
			if path == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			u.Response(w, http.StatusUnauthorized, u.Message(false, "Malformed auth token"))
			return
		}

		splicedTokenHeader := strings.Split(tokenHeader, " ")
		if len(splicedTokenHeader) != 2 {
			u.Response(w, http.StatusUnauthorized, u.Message(false, "Invalid/Malformed auth token"))
			return
		}

		tokenPart := splicedTokenHeader[1]
		tk := &models.Token{}
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
		if err != nil {
			u.Response(w, http.StatusUnauthorized, u.Message(false, "Malformed auth token"))
			return
		}

		var cnt int
		models.GetDB().Table("tokens_blacklist").Where("token = ?", token.Raw).Count(&cnt)

		if !token.Valid || cnt == 1 {
			u.Response(w, http.StatusUnauthorized, u.Message(false, "Token is not valid"))
			return
		}
		if tk.ExpirationTime < time.Now().Unix() {
			u.Response(w, http.StatusUnauthorized, u.Message(false, "Token expired"))
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
