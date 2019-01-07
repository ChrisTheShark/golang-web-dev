package main

import (
	"log"
	"net/http"
	"text/template"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFiles("templates/subscribe.gohtml"))
}

type Person struct {
	Firstname  string
	Lastname   string
	Subscribed bool
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var person Person

		person.Firstname = req.FormValue("first_name")
		person.Lastname = req.FormValue("last_name")

		value := req.FormValue("subscribe")
		if value == "on" {
			person.Subscribed = true
		}

		log.Println(person)

		templates.ExecuteTemplate(w, "subscribe.gohtml", person)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
