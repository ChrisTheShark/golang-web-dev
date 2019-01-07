package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Friend struct {
	ID   int
	Name string
}

var db *sql.DB

func init() {
	var err error
	db, err = getDatabase()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS friends (
			friend_id INT AUTO_INCREMENT, 
			name VARCHAR(255) NOT NULL,
			PRIMARY KEY (friend_id)
		)
	`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`DELETE FROM friends`); err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer db.Close()

	http.HandleFunc("/friends/", func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/friends/"))
		if err != nil {
			log.Println(err)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		if r.Method == http.MethodGet {
			friend, err := getByID(id)
			if err != nil {
				log.Println(err)
				if err == sql.ErrNoRows {
					log.Println(err)
					http.Error(w, "not found", http.StatusNotFound)
					return
				}

				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			}

			res, err := json.Marshal(friend)
			if err != nil {
				log.Println(err)
				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			}

			w.Write(res)
			return
		}

		if r.Method == http.MethodDelete {
			if err := deleteByID(id); err != nil {
				log.Println(err)
				if err == sql.ErrNoRows {
					log.Println(err)
					http.Error(w, "not found", http.StatusNotFound)
					return
				}

				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			}

			w.WriteHeader(204)
			return
		}

	})

	http.HandleFunc("/friends", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			friends, err := getAll()
			if err != nil {
				log.Println(err)
				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			}

			if len(friends) > 0 {
				res, err := json.Marshal(friends)
				if err != nil {
					log.Println(err)
					http.Error(w, "service unavailable", http.StatusServiceUnavailable)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(res)
				return
			}

			w.WriteHeader(204)
			return
		}

		if r.Method == http.MethodPost {
			var friend Friend

			if bytes, err := ioutil.ReadAll(r.Body); err != nil {
				log.Println(err)
				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			} else {
				if err := json.Unmarshal(bytes, &friend); err != nil {
					log.Println(err)
					http.Error(w, "service unavailable", http.StatusServiceUnavailable)
					return
				}
			}

			f, err := add(friend.Name)
			if err != nil {
				log.Println(err)
				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/friends/%d", f.ID), http.StatusSeeOther)
			return
		}

	})

	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getDatabase() (*sql.DB, error) {
	// "user:password@tcp(host:port)/database?charset=utf8"
	db, err := sql.Open("mysql",
		"root:password@tcp(localhost:3306)/sample?charset=utf8")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
