package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/index.css" {
		http.ServeFile(w, r, "template/index.css")
		return
	}

	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		text := r.FormValue("text")
		banner := r.FormValue("banner")
		if banner == "" {
			banner = "standard"
		}
		result, err := ascii(text, banner)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, result)
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", Handler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}