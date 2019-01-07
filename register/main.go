package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username  string
	Password  []byte
	Firstname string
	Lastname  string
}

func (u user) String() string {
	return fmt.Sprintf(`
		Username: %s,
		Password: %s,
		Firstname: %s,
		Lastname: %s,
	`, u.Username, string(u.Password), u.Firstname, u.Lastname)
}

var templates *template.Template
var sessions = map[string]string{}
var users = map[string]user{}

func init() {
	templates = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var user user
		cookie, err := r.Cookie("session")
		if err != nil {
			log.Printf("error locating session - %s\n", err)
		} else {
			username := sessions[cookie.Value]
			user = users[username]
		}
		templates.ExecuteTemplate(w, "index.gohtml", user)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if hasActiveSession(r) {
			log.Println("user has active session, redirecting.")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		templates.ExecuteTemplate(w, "login.gohtml", nil)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if !hasActiveSession(r) {
			log.Println("user has no active session, redirecting.")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		cookie, err := r.Cookie("session")
		if err != nil {
			log.Println(err)
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		cookie.MaxAge = -1
		http.SetCookie(w, cookie)

		delete(sessions, cookie.Value)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			encrypted, err := bcrypt.GenerateFromPassword(
				[]byte(r.FormValue("password")), bcrypt.DefaultCost)
			if err != nil {
				log.Println(err)
				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			}

			user := user{
				Username:  r.FormValue("user_name"),
				Password:  encrypted,
				Firstname: r.FormValue("first_name"),
				Lastname:  r.FormValue("last_name"),
			}
			users[r.FormValue("user_name")] = user

			cookie, err := createSession()
			if err != nil {
				log.Println(err)
				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			}

			log.Printf("created new user: %s and linked to session: %s\n",
				user, cookie.Value)

			sessions[cookie.Value] = user.Username
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username := r.FormValue("user_name")
			user, ok := users[username]
			if !ok {
				log.Printf("Unable to locate user with username: %s", username)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			if err := bcrypt.CompareHashAndPassword(user.Password,
				[]byte(r.FormValue("password"))); err != nil {
				log.Printf("invalid credentials provided for username: %s", username)
				http.Error(w, "username or password invalid", http.StatusUnauthorized)
				return
			}

			cookie, err := createSession()
			if err != nil {
				log.Println(err)
				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			}

			log.Printf("logged in user: %s and linked to session: %s\n",
				user, cookie.Value)

			sessions[cookie.Value] = user.Username
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
