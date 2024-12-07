package main

import (
	// "fmt"
	"go-web-native/config"
	"go-web-native/config/middleware"
	"go-web-native/controllers/authcontroller"
	"go-web-native/controllers/categorycontroller"
	"go-web-native/controllers/homecontroller"
	"go-web-native/controllers/productcontroller"
	"go-web-native/templates"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func init()  {
	templates.Init()
}

// Fungsi untuk melayani file statis dengan pengecekan apakah path adalah file
func staticFileHandler(w http.ResponseWriter, r *http.Request) {
    filePath := filepath.Join("public", r.URL.Path[len("/public/"):])

    // Cek apakah path mengarah ke file, bukan folder
    fileInfo, err := os.Stat(filePath)
    if os.IsNotExist(err) || fileInfo.IsDir() {
        http.NotFound(w, r) // Jika file tidak ada atau adalah folder, kembalikan 404
        return
    }

    // Jika path valid dan mengarah ke file, layani file
    http.ServeFile(w, r, filePath)
}

func main() {
	config.ConnectDB()
	routes()

	http.HandleFunc("/public/", staticFileHandler)
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func routes()  {
	// 1. Authentication
	http.HandleFunc("/login", authcontroller.Login)
	http.HandleFunc("/register", authcontroller.Register)
	http.HandleFunc("/logout", authcontroller.Logout)

	// 2. Homepage
	http.HandleFunc("/", homecontroller.Welcome)

	// 3. Categories
	http.HandleFunc("/categories", middleware.IsLogin(categorycontroller.Index))
	http.HandleFunc("/categories/add", middleware.IsLogin(categorycontroller.Add))
	http.HandleFunc("/categories/edit", middleware.IsLogin(categorycontroller.Edit))
	http.HandleFunc("/categories/delete", middleware.IsLogin(categorycontroller.Delete))

	// 4. Products
	http.HandleFunc("/products", middleware.IsLogin(productcontroller.Index))
	http.HandleFunc("/products/add", middleware.IsLogin(productcontroller.Add))
	http.HandleFunc("/products/edit", middleware.IsLogin(productcontroller.Edit))
	http.HandleFunc("/products/delete", middleware.IsLogin(productcontroller.Delete))
}