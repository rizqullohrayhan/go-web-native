package middleware

import (
	"net/http"

	"github.com/kataras/go-sessions"
)

// Middleware untuk memeriksa apakah pengguna sudah login
func IsLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sessions.Start(w, r)
		if session.GetString("username") == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}