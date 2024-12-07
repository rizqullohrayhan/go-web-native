package auth

import (
	"go-web-native/entities"
	"go-web-native/models/usermodel"
	"net/http"

	"github.com/kataras/go-sessions"
	"golang.org/x/crypto/bcrypt"
)

func Attempt(w http.ResponseWriter, r *http.Request, username string, password string) bool {
	userFromDB, err := usermodel.GetByUsername(username)
	if err != nil || userFromDB == nil {
		return false
	}
	checkPass := bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(password))
	if checkPass != nil {
		return false
	}
	setSession(w, r, *userFromDB)
	return true
}

func Login(w http.ResponseWriter, r *http.Request, username string) bool {
	userFromDB, err := usermodel.GetByUsername(username)
	if err != nil || userFromDB == nil {
		return false
	}
	setSession(w, r, *userFromDB)
	return true
}

func User(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	session := sessions.Start(w, r)
	if session.GetString("id") != "" {
		user := session.GetAll()
		return user
	}
	return nil
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Destroy()
}

func setSession(w http.ResponseWriter, r *http.Request, user entities.User)  {
	session := sessions.Start(w, r)
	session.Set("id", user.ID)
	session.Set("username", user.Username)
	session.Set("name", user.FirstName)
	// session.Set("isLoggedIn", true)
}