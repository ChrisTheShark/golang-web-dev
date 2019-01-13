package main

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("./templates/*.gohtml"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := getCookie(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		xs := strings.Split(cookie.Value, "|")
		templates.ExecuteTemplate(w, "index.gohtml", xs[1:])
	})
	http.HandleFunc("/photos", func(w http.ResponseWriter, r *http.Request) {
		f, fh, err := r.FormFile("photo")
		if err != nil {
			log.Println(err)
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer f.Close()

		body, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		h := sha1.New()
		ext := strings.Split(fh.Filename, ".")[1]

		io.Copy(h, f)

		filename := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext

		wd, err := os.Getwd()
		if err != nil {
			log.Println(err)
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		err = ioutil.WriteFile(filepath.Join(wd, "public", "images", filename),
			body, os.ModePerm)
		if err != nil {
			log.Println(err)
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		cookie, err := getCookie(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		appendValue(w, cookie, filename)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	http.Handle("/public/", http.FileServer(http.Dir(".")))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getCookie(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		sessionID, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}

		cookie = &http.Cookie{
			Name:  "session",
			Value: sessionID.String(),
		}
		http.SetCookie(w, cookie)
	}
	return cookie, nil
}

func appendValue(w http.ResponseWriter, cookie *http.Cookie, file string) *http.Cookie {
	if !strings.Contains(cookie.Value, file) {
		cookie.Value = cookie.Value + "|" + file
	}
	http.SetCookie(w, cookie)
	return cookie
}
