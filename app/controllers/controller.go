package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var Tmpl *template.Template

// Create a home handler:
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
	err := Tmpl.ExecuteTemplate(w, "home.html", "Home page!!!")
	if err != nil {
		log.Fatal(err)
	}
}
