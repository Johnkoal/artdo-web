package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func main2() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", home)

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Host", r.Host)
		fmt.Fprintln(w, "RequestURI", r.RequestURI)
		fmt.Fprintln(w, "Method", r.Method)
		fmt.Fprintln(w, "RemoteAddr", r.RemoteAddr)
	})

	http.HandleFunc("/product", product)

	http.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/product", 301)
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Has error generated", 404)
	})

	http.HandleFunc("/head", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Test", "test1")

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintln(w, "{ \"hola\":1 }")
	})

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, struct{ Saludo string }{"Holis"})
	})

	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	html := "<html>"
	html += "<body>"
	html += "<h1>Hello, World</h1>"
	html += "</body>"
	html += "</html>"
	w.Write([]byte(html))
}

var products []string

func product(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	add, okForm := r.Form["add"]
	if okForm && len(add) == 1 {
		products = append(products, string(add[0]))
		w.Write([]byte("Product add successful"))

		return
	}

	prod, ok := r.URL.Query()["produc"]
	if ok && len(prod) == 1 {
		pos, err := strconv.Atoi(prod[0])

		if err != nil {
			return
		}

		html := "<html>"
		html += "<body>"
		html += "<h1>The product selected: " + products[pos] + "</h1>"
		html += "</body>"
		html += "</html>"
		w.Write([]byte(html))

		return
	}

	html := "<html>"
	html += "<body>"
	html += "<h1>Total Produts " + strconv.Itoa(len(products)) + "</h1>"
	html += "</body>"
	html += "</html>"
	w.Write([]byte(html))

}
