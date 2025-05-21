package login

import (
	"net/http"

	"github.com/kirill-dolgii/url-shortner/internal/clients/sso/grpc"
)

func New(client *grpc.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")
		appID := 1 // например, идентификатор приложения, откуда пользователь заходит

		token, err := client.Login(r.Context(), email, password, appID)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
		})

		w.Write([]byte("login successful"))
	}
}
