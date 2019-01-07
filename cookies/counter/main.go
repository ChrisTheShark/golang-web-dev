package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("counter")
		if err == http.ErrNoCookie {
			http.SetCookie(w, &http.Cookie{
				Name:  "counter",
				Value: "1",
			})
			io.WriteString(w, "Welcome to the site!")
			return
		}

		num, err := strconv.Atoi(cookie.Value)
		log.Printf("Located visit number: %v", num)

		http.SetCookie(w, &http.Cookie{
			Name:  "counter",
			Value: strconv.Itoa(num + 1),
		})
		io.WriteString(w,
			fmt.Sprintf("Welcome back, I see this is your visit number: %v!", num+1))
	})
	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("counter")
		if err == http.ErrNoCookie {
			log.Println(err)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		// Set MaxAge to negative value removes the cookie.
		cookie.MaxAge = -1

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
