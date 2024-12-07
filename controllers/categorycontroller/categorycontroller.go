package categorycontroller

import (
	"fmt"
	"go-web-native/models/categorymodel"
	"go-web-native/templates"
	"net/http"

	"github.com/kataras/go-sessions"
)

func checkSession(w http.ResponseWriter, r *http.Request) bool {
	session := sessions.Start(w, r)
	isLoggedIn := session.GetString("username") != ""
	return isLoggedIn
}

func Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{} {
		"categories": categorymodel.GetAll(),
		"isLoggedIn": checkSession(w, r),
	}

	err := templates.CategoryTemplates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{} {
		"isLoggedIn": checkSession(w, r),
	}
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
		name := r.FormValue("name")
		if name == "" {
			err := templates.CategoryTemplates.ExecuteTemplate(w, "create.html", data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if err := categorymodel.Create(name); err != nil {
			err := templates.CategoryTemplates.ExecuteTemplate(w, "create.html", data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
	case "GET":
		err := templates.CategoryTemplates.ExecuteTemplate(w, "create.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{} {
		"category": categorymodel.GetByID(r.URL.Query().Get("id")),
		"isLoggedIn": checkSession(w, r),
	}
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
		name := r.FormValue("name")
		id := r.URL.Query().Get("id")
		if name == "" {
			err := templates.CategoryTemplates.ExecuteTemplate(w, "edit.html", data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if err := categorymodel.Edit(name, id); err != nil {
			err := templates.CategoryTemplates.ExecuteTemplate(w, "edit.html", data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
	case "GET":
		err := templates.CategoryTemplates.ExecuteTemplate(w, "edit.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

	id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "Missing category ID", http.StatusBadRequest)
        return
    }

	if err := categorymodel.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
        return
	}
	w.WriteHeader(http.StatusOK)
}