package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "Session_ID",
			Value: "1234567890",
		})
		http.Redirect(w, req, "/get", http.StatusSeeOther)
	})
	http.HandleFunc("/get", func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("Session_ID")
		if err != nil {
			log.Println(err)
			http.Error(w, "not found", http.StatusNotFound)
		}
		io.WriteString(w, cookie.String())
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
