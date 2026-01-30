package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
)

type PageData struct {
	Title           string
	MetaDescription string            // Descripción SEO traducida
	MetaKeywords    string            // Keywords SEO
	CanonicalURL    string            // URL canónica
	T               map[string]string // Mapa de traducción
	Lang            string            // Idioma actual ("es" o "en")
	Config          Config            // Configuración global
	SuccessMessage  string            // Mensaje de éxito
	ErrorMessage    string            // Mensaje de error
}

type Config struct {
	CompanyName string `json:"company_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Social      struct {
		Facebook  string `json:"facebook"`
		Twitter   string `json:"twitter"`
		Instagram string `json:"instagram"`
		Linkedin  string `json:"linkedin"`
	} `json:"social"`
	SMTP struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"smtp"`
	Recaptcha struct {
		SiteKey   string `json:"site_key"`
		SecretKey string `json:"secret_key"`
	} `json:"recaptcha"`
}

var translations = make(map[string]map[string]string)
var globalConfig Config

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error abriendo config.json: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&globalConfig)
	if err != nil {
		log.Fatalf("Error decodificando config.json: %v", err)
	}
	log.Println("Configuración global cargada.")
}

func loadTranslations() {
	files := []string{"es", "en"}
	for _, lang := range files {
		file, err := os.Open(fmt.Sprintf("locales/%s.json", lang))
		if err != nil {
			log.Fatalf("Error abriendo archivo de idioma %s: %v", lang, err)
		}
		defer file.Close()

		byteValue, _ := ioutil.ReadAll(file)
		var result map[string]string
		json.Unmarshal(byteValue, &result)
		translations[lang] = result
		log.Printf("Idioma cargado: %s", lang)
	}
}

func main() {
	loadConfig()
	loadTranslations()

	// Servir archivos estáticos
	http.Handle("/artdotech-core.css", http.FileServer(http.Dir(".")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handleRequest("index.html", "home"))
	http.HandleFunc("/about", handleRequest("about.html", "about"))
	// Mantener compatibilidad con rutas antiguas (redirección)
	http.HandleFunc("/quienes-somos", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/about"+langQuery(r), http.StatusMovedPermanently)
	})

	http.HandleFunc("/services", handleRequest("services.html", "services"))
	// Redirección antigua
	http.HandleFunc("/nuestros-servicios", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/services"+langQuery(r), http.StatusMovedPermanently)
	})

	http.HandleFunc("/contact", handleContact)

	// Redirección antigua
	http.HandleFunc("/contacto", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/contact"+langQuery(r), http.StatusMovedPermanently)
	})

	// SEO Files
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/robots.txt")
	})
	http.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/sitemap.xml")
	})

	log.Println("Servidor escuchando en :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Helper para obtener datos comunes de la página
func getPageData(r *http.Request, pageType string) PageData {
	lang := r.URL.Query().Get("lang")
	if lang == "" || (lang != "es" && lang != "en") {
		lang = "es" // Default
	}

	// Títulos base también traducidos
	titles := map[string]string{
		"home":     translations[lang]["meta_title_home"],
		"about":    translations[lang]["meta_title_about"],
		"services": translations[lang]["meta_title_services"],
		"contact":  translations[lang]["meta_title_contact"],
	}

	// Meta descripciones
	descriptions := map[string]string{
		"home":     translations[lang]["meta_description_home"],
		"about":    translations[lang]["meta_description_about"],
		"services": translations[lang]["meta_description_services"],
		"contact":  translations[lang]["meta_description_contact"],
	}

	// Construir Canonical URL (Siempre apunta al dominio principal sin query params de idioma si es el default, o con ellos si es explícito)
	// Para simplificar, apuntaremos a la versión oficial de esa página
	path := r.URL.Path
	if path == "/" {
		path = ""
	}
	canonical := "https://artdotech.com" + path
	if lang == "en" {
		canonical += "?lang=en"
	}

	return PageData{
		Title:           titles[pageType],
		MetaDescription: descriptions[pageType],
		MetaKeywords:    translations[lang]["meta_keywords"],
		CanonicalURL:    canonical,
		T:               translations[lang],
		Lang:            lang,
		Config:          globalConfig,
	}
}

// Manejador genérico para páginas estáticas
func handleRequest(tmpl string, pageType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := getPageData(r, pageType)
		renderTemplate(w, tmpl, data)
	}
}

// Manejador específico para el formulario de contacto
func handleContact(w http.ResponseWriter, r *http.Request) {
	data := getPageData(r, "contact")

	if r.Method == http.MethodPost {
		// Procesar el formulario
		err := r.ParseForm()
		if err != nil {
			data.ErrorMessage = "Error procesando el formulario."
			renderTemplate(w, "contact.html", data)
			return
		}

		// Validar CAPTCHA
		recaptchaResponse := r.FormValue("g-recaptcha-response")
		if !verifyCaptcha(recaptchaResponse) {
			data.ErrorMessage = "Por favor verifica que no eres un robot."
			renderTemplate(w, "contact.html", data)
			return
		}

		nombre := r.FormValue("nombre")
		email := r.FormValue("email")
		asunto := r.FormValue("asunto")
		mensaje := r.FormValue("mensaje")

		// Enviar correo
		err = sendEmail(nombre, email, asunto, mensaje)
		if err != nil {
			log.Printf("Error enviando correo: %v", err)
			data.ErrorMessage = "Hubo un problema enviando tu mensaje. Por favor intenta más tarde."
		} else {
			data.SuccessMessage = "¡Mensaje enviado con éxito! Nos pondremos en contacto contigo pronto."
		}
	}

	renderTemplate(w, "contact.html", data)
}

func sendEmail(nombre, email, asunto, mensaje string) error {
	smtpConfig := globalConfig.SMTP
	// Si no hay configuración SMTP, solo loguear (modo desarrollo/mock)
	if smtpConfig.Host == "" {
		log.Printf("[MOCK EMAIL] De: %s <%s> | Asunto: %s | Mensaje: %s", nombre, email, asunto, mensaje)
		return nil
	}

	from := smtpConfig.Username
	password := smtpConfig.Password
	to := globalConfig.Email
	smtpHost := smtpConfig.Host
	smtpPort := fmt.Sprintf("%d", smtpConfig.Port)

	// Mensaje formateado
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Nuevo mensaje de contacto web: %s\r\n"+
		"\r\n"+
		"Nombre: %s\r\n"+
		"Email: %s\r\n"+
		"Asunto: %s\r\n"+
		"\r\n"+
		"Mensaje:\r\n%s\r\n", to, asunto, nombre, email, asunto, mensaje))

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	return err
}

func verifyCaptcha(token string) bool {
	secretKey := globalConfig.Recaptcha.SecretKey
	if secretKey == "" || secretKey == "TU_CLAVE_SECRETA" {
		// En desarrollo, si no hay clave, pasa (o podrías bloquearlo)
		log.Println("[MOCK CAPTCHA] Verificación saltada (clave no configurada).")
		return true
	}

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {secretKey}, "response": {token}})
	if err != nil {
		log.Printf("Error conectando con reCAPTCHA: %v", err)
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decodificando respuesta reCAPTCHA: %v", err)
		return false
	}

	return result.Success
}

// Helper para mantener el query param de idioma en redirecciones
func langQuery(r *http.Request) string {
	lang := r.URL.Query().Get("lang")
	if lang != "" {
		return "?lang=" + lang
	}
	return ""
}

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
