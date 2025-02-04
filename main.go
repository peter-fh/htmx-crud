package main

import (
	"net/http"
	"html/template"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var tmpl *template.Template


func init() {
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	check(err)
}

func Homepage (w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.html", nil)
}

func main() {
	http.HandleFunc("/", Homepage)
	http.ListenAndServe(":4000", nil)
}
