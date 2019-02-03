package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestGetAllPeople(t *testing.T) {
	people = map[int]Person{
		1: Person{
			Name: "Jon Thomas",
			Age:  32,
			Friends: []string{
				"Thomas",
				"Reign",
			},
		},
	}

	r := httptest.NewRequest(http.MethodGet, "/people", nil)
	w := httptest.NewRecorder()

	getAllPeople(w, r)
	resp := w.Result()

	bs, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, "[{\"name\":\"Jon Thomas\",\"age\":32,\"friends\":[\"Thomas\",\"Reign\"]}]\n", string(bs))
}

func TestGetPersonByID(t *testing.T) {
	people = map[int]Person{
		1: Person{
			Name: "Jon Thomas",
			Age:  32,
			Friends: []string{
				"Thomas",
				"Reign",
			},
		},
	}

	r := httptest.NewRequest(http.MethodGet, "/people/1", nil)
	w := httptest.NewRecorder()

	getPersonByID(w, r)
	resp := w.Result()

	bs, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, "{\"name\":\"Jon Thomas\",\"age\":32,\"friends\":[\"Thomas\",\"Reign\"]}\n", string(bs))
}

func TestGetPersonByIDNotFound(t *testing.T) {
	people = map[int]Person{}

	r := httptest.NewRequest(http.MethodGet, "/people/1", nil)
	w := httptest.NewRecorder()

	getPersonByID(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "404 Not Found", resp.Status)
}

func TestAddPerson(t *testing.T) {
	people = map[int]Person{}
	person := Person{
		Name: "Amy Jonson",
		Age:  39,
		Friends: []string{
			"Jak",
			"Daxter",
		},
	}

	bs, _ := json.Marshal(&person)
	r := httptest.NewRequest(http.MethodPost, "/people", bytes.NewReader(bs))
	w := httptest.NewRecorder()

	addPerson(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
	assert.Equal(t, "303 See Other", resp.Status)
	assert.Equal(t, people[1], person)
}

func TestAddPersonBadRequest(t *testing.T) {
	otherStruct := struct {
		Sport      string
		NumPlayers int
	}{
		Sport:      "basketball",
		NumPlayers: 5,
	}

	bs, _ := json.Marshal(&otherStruct)
	r := httptest.NewRequest(http.MethodPost, "/people", bytes.NewReader(bs))
	w := httptest.NewRecorder()

	addPerson(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "400 Bad Request", resp.Status)
}
