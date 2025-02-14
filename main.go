package main

import (
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Title string
}

func main() {
	// Servir archivos estáticos (CSS, JS, imágenes, etc.)
	http.Handle("/artdotech-core.css", http.FileServer(http.Dir(".")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", inicio) // Ruta para la página de inicio
	http.HandleFunc("/quienes-somos", quienesSomos)
	http.HandleFunc("/nuestros-servicios", nuestrosServicios)
	http.HandleFunc("/contacto", contacto)

	// Redirigir la raíz a "quienes-somos"
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// })

	// Iniciar el servidor
	log.Println("Servidor escuchando en :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func inicio(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Inicio - Artdotech"}
	renderTemplate(w, "index.html", data)
}

func quienesSomos(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "¿Quiénes somos? - Artdotech"}
	renderTemplate(w, "quienes_somos.html", data)
}

func nuestrosServicios(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Nuestros servicios - Artdotech"}
	renderTemplate(w, "nuestros_servicios.html", data)
}

func contacto(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Contáctenos - Artdotech"}
	renderTemplate(w, "contacto.html", data)
}

// Función para renderizar templates
func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	t, err := template.ParseFiles("templates/layout.html", "templates/"+tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
