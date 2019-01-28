package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/get", get)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "Session_ID",
		Value: "1234567890",
	})
	http.Redirect(w, req, "/get", http.StatusSeeOther)
}

func get(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("Session_ID")
	if err != nil {
		log.Println(err)
		http.Error(w, "not found", http.StatusNotFound)
	}
	io.WriteString(w, cookie.String())
}
