package login

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			tokenStr := cookie.Value
			claims := jwt.MapClaims{}
			tn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			fmt.Println(tn.Method)

			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			userID := int(claims["uid"].(float64))
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
