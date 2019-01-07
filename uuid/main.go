package main

import (
	"io"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			id, err := uuid.NewV4()
			if err != nil {
				log.Println(err)
				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			}

			cookie = &http.Cookie{
				Name:     "session",
				Value:    id.String(),
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}

		io.WriteString(w, cookie.String())
	})
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			log.Println(err)
			http.Error(w, "not found", http.StatusNotFound)
		}

		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	})
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
