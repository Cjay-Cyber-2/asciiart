package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func  main(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		if r.URL.Path != "/"{
		http.Error(w, "404 not found", http.StatusNotFound)
		return
		}

		tmpl, err := template.ParseFiles("index.html")
		if err != nil{
			http.Error(w,"500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, nil)
	})


	http.

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}