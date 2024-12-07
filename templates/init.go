package templates

import (
	"html/template"
	"path/filepath"
)

var (
	HomeTemplates     *template.Template
	ProductTemplates  *template.Template
	CategoryTemplates *template.Template
	AuthTemplates *template.Template
)

// loadTemplates adalah helper function untuk menginisialisasi template dengan layout
func loadTemplates(folder string) *template.Template {
	// Mulai dengan mem-parsing file layout
	tmpl := template.Must(template.ParseGlob("views/layout/*.html"))

	// Gabungkan dengan file template dari folder spesifik
	template.Must(tmpl.ParseGlob(filepath.Join("views", folder, "*.html")))

	return tmpl
}

func Init() {
	AuthTemplates = loadTemplates("auth")
	HomeTemplates = loadTemplates("home")
	ProductTemplates = loadTemplates("product")
	CategoryTemplates = loadTemplates("category")
}