package authcontroller

import (
	"go-web-native/config/auth"
	"go-web-native/entities"
	"go-web-native/models/usermodel"
	"go-web-native/templates"
	"log"
	"net/http"

	// "github.com/kataras/go-sessions"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.AuthTemplates.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var user entities.User
	user.Username = r.FormValue("email")
	user.FirstName = r.FormValue("first_name")
	user.LastName = r.FormValue("last_name")
	user.Password = r.FormValue("password")

	// Cek apakah username sudah ada
	userFromDB, _ := usermodel.GetByUsername(user.Username)
	if userFromDB != nil {
		http.Error(w, "Username sudah digunakan", http.StatusBadRequest)
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Gagal memproses password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword) // Simpan password yang sudah di-hash

	// Simpan pengguna baru ke database
	err = usermodel.Create(user)
	if err != nil {
		http.Error(w, "Gagal menyimpan pengguna", http.StatusInternalServerError)
		return
	}

	// Redirect atau tampilkan pesan berhasil
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// log.Println(auth.User(w, r))
	if auth.User(w, r) != nil {
		log.Println("logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		log.Println("not logged in")
		err := templates.AuthTemplates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	username := r.FormValue("email")
	password := r.FormValue("password")


	// userFromDB, err := usermodel.GetByUsername(username)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// if userFromDB == nil {
	// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
	// 	return
	// }
	// password_tes := bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(password))

	if auth.Attempt(w, r, username, password) {
		//login success
		// session.Set("username", userFromDB.Username)
		// session.Set("name", userFromDB.FirstName)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		//login failed
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}