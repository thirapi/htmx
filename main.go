package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Current time: %s", time.Now().Format(time.RFC1123))
	})

	fmt.Println("Server starting at :8080...")
	http.ListenAndServe(":8080", nil)
}
