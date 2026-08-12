package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var allowedBanners = map[string]bool{
	"standard":   true,
	"shadow":     true,
	"thinkertoy": true,
}

type pageData struct {
	Result string
	Text   string
	Banner string
	Error  string
}

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

	data := pageData{Banner: "standard"}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.Error = "Invalid form data."
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, data)
			return
		}

		text := r.FormValue("text")
		banner := r.FormValue("banner")
		if banner == "" {
			banner = "standard"
		}
		data.Text = text
		data.Banner = banner

		if !allowedBanners[banner] {
			data.Error = "Unknown banner. Choose standard, shadow, or thinkertoy."
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, data)
			return
		}

		if strings.TrimSpace(text) == "" {
			data.Error = "Please enter some text to convert."
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, data)
			return
		}

		result, err := ascii(text, banner)
		if err != nil {
			data.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			tmpl.Execute(w, data)
			return
		}

		data.Result = result
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func ascii(input, banner string) (string, error) {
	if !allowedBanners[banner] {
		return "", fmt.Errorf("unknown banner")
	}

	content, err := os.ReadFile("banners/" + banner + ".txt")
	if err != nil {
		return "", fmt.Errorf("could not load banner file")
	}

	normalized := strings.ReplaceAll(string(content), "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	words := strings.Split(normalized, "\n")

	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")
	input = strings.ReplaceAll(input, "\\n", "\n")
	lines := strings.Split(input, "\n")

	var result strings.Builder

	for _, line := range lines {
		if line == "" {
			result.WriteString("\n")
			continue
		}

		for row := 0; row < 8; row++ {
			for _, c := range line {
				code := int(c)
				if code < 32 || code > 126 {
					return "", fmt.Errorf("input contains unsupported characters (use printable ASCII only)")
				}
				index := (code-32)*9 + 1 + row
				if index < 0 || index >= len(words) {
					return "", fmt.Errorf("banner file is invalid")
				}
				result.WriteString(words[index])
			}
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

func main() {
	http.HandleFunc("/", Handler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
