package homecontroller

import (
	"go-web-native/config/auth"
	"go-web-native/templates"
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"auth": auth.User(w, r),
	}
	err := templates.HomeTemplates.ExecuteTemplate(w, "index.html", data) // "index.html" sesuai dengan nama file di folder template
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}