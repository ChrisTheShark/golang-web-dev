package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Finally, I've been expecting you!!")
	})
	http.HandleFunc("/foo", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusSeeOther)
	})
	http.HandleFunc("/bar", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
