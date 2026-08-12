package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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

func ascii(input, banner string) (string, error) {
	filename := "banners/" + banner + ".txt"
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	words := strings.Split(string(content), "\n")
	lines := strings.Split(input, "\\n")

	var result strings.Builder

	for _, char := range lines {
		if char == "" {
			result.WriteString("\n")
			continue
		}
		for row := 0; row < 9; row++ {
			for _, c := range char {
				code := int(c)
				if code >= 32 && code <= 126 {
					index := (code - 32) * 9 + 1
					result.WriteString(words[index+row])
				}
			}
			if row < 8 {
				result.WriteString("\n")
			}
		}
	}
	return result.String(), nil
}

func main() {
	http.HandleFunc("/", Handler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}