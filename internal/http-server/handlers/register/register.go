package register

import (
	"net/http"

	"github.com/kirill-dolgii/url-shortner/internal/clients/sso/grpc"
)

func New(client *grpc.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		_, err := client.Register(r.Context(), email, password)
		if err != nil {
			http.Error(w, "registration failed", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("registration successful"))
	}
}
