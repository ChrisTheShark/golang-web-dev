package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"encoding/json"
)

func main() {
	http.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getAllPeople(w, r)
			return
		}
		if r.Method == http.MethodPost {
			addPerson(w, r)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})

	http.HandleFunc("/people/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getPersonByID(w, r)
			return
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAllPeople(w http.ResponseWriter, r *http.Request) {
	ps, err := getAll()
	if err != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ps)
	return
}

func getPersonByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	person, err := getByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

func addPerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil || person.isEmpty() {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	id, err := add(person)
	if err != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/people/%v", id), http.StatusSeeOther)
}

func parseID(r *http.Request) (int, error) {
	split := strings.Split(r.URL.String(), "/")
	if len(split) > 2 {
		id, err := strconv.Atoi(split[2])
		if err != nil {
			return 0, errors.New("unable to parse provided ID")
		}
		return id, nil
	}
	return 0, errors.New("no id located in URL")
}
