package productcontroller

import (
	"fmt"
	"go-web-native/entities"
	"go-web-native/models/categorymodel"
	"go-web-native/models/productmodel"
	"go-web-native/templates"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kataras/go-sessions"
	// "golang.org/x/telemetry/upload"
)

func checkSession(w http.ResponseWriter, r *http.Request) bool {
	session := sessions.Start(w, r)
	isLoggedIn := session.GetString("username") != ""
	return isLoggedIn
}

func Index(w http.ResponseWriter, r *http.Request) {
	products := productmodel.GetAll()
	data := map[string]any {
		"products": products,
		"isLoggedIn": checkSession(w, r),
	}

	err := templates.ProductTemplates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	categories := categorymodel.GetAll()
	data := map[string]any {
		"categories": categories,
		"isLoggedIn": checkSession(w, r),
	}

	switch r.Method {
	case "POST":
		if err := r.ParseMultipartForm(1024); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var product entities.Product
		product.Name = r.FormValue("name")
		categoryId, _ := strconv.ParseUint(r.FormValue("category_id"), 10, 64)
		product.Stock, _ = strconv.Atoi(r.FormValue("stock"))
		product.CategoryId = uint(categoryId)
		product.Description = r.FormValue("description")
		if product.Name == "" || product.Description == "" || product.Stock < 0 {
			viewCreate(w, data)
			return
		}

		// Handling Uploaded File
		uploadedFile, handler, err := r.FormFile("image")
		switch err {
		case nil:
			dir, err := os.Getwd()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			filename := handler.Filename
			fileLocation := filepath.Join(dir, "public/upload/product", filename)
			targetFile, err := os.OpenFile(fileLocation, os.O_RDONLY|os.O_CREATE, 0666)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer targetFile.Close()

			if _, err := io.Copy(targetFile, uploadedFile); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			product.Image = filename
		case http.ErrMissingFile:
			viewCreate(w, data)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer uploadedFile.Close()

		if err := productmodel.Create(product); err != nil {
			log.Printf("Error: %v", err.Error())
			viewCreate(w, data)
			return
		}
		http.Redirect(w, r, "/products", http.StatusSeeOther)
	case "GET":
		viewCreate(w, data)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	Id := r.URL.Query().Get("id")
	categories := categorymodel.GetAll()
	data := map[string]any {
		"categories": categories,
		"product": productmodel.GetById(Id),
		"isLoggedIn": checkSession(w, r),
	}

	switch r.Method {
	case "POST":
		if err := r.ParseMultipartForm(1024); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
		var product entities.Product
		product.Name = r.FormValue("name")
		categoryId, _ := strconv.ParseUint(r.FormValue("category_id"), 10, 64)
		product.Stock, _ = strconv.Atoi(r.FormValue("stock"))
		product.CategoryId = uint(categoryId)
		product.Description = r.FormValue("description")
		if product.Name == "" || product.Description == "" || product.Stock < 0 {
			viewEdit(w, data)
			return
		}

		// Handling Uploaded File
		uploadedFile, handler, err := r.FormFile("image")
		oldImage := productmodel.GetById(Id).Image
		switch err {
		case nil:
			defer uploadedFile.Close()
			dir, err := os.Getwd()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			filename := handler.Filename
			fileLocation := filepath.Join(dir, "public/upload/product", filename)
			targetFile, err := os.OpenFile(fileLocation, os.O_RDONLY|os.O_CREATE, 0666)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer targetFile.Close()

			// Hapus file lama
			if oldImage != "" {
				oldImagePath := filepath.Join(dir, "public/upload/product", oldImage)
				err = os.Remove(oldImagePath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			if _, err := io.Copy(targetFile, uploadedFile); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			product.Image = filename
		case http.ErrMissingFile:
			// jika user tidak upload file
			product.Image = oldImage
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := productmodel.Edit(product, Id); err != nil {
			viewEdit(w, data)
			return
		}
		http.Redirect(w, r, "/products", http.StatusSeeOther)
	case "GET":
		viewEdit(w, data)
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

	if err := productmodel.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
        return
	}
	w.WriteHeader(http.StatusOK)
}

func viewCreate(w http.ResponseWriter, data any) {
	err := templates.ProductTemplates.ExecuteTemplate(w, "create.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewEdit(w http.ResponseWriter, data any) {
	err := templates.ProductTemplates.ExecuteTemplate(w, "edit.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}