package main

import "errors"

// Person type represents a person object.
type Person struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

func (p Person) isEmpty() bool {
	return p.Name == "" && p.Age == 0 && p.Friends == nil
}

// Simple data structure to hold people.
var people = map[int]Person{}

func getAll() ([]Person, error) {
	ps := []Person{}
	for _, value := range people {
		ps = append(ps, value)
	}
	return ps, nil
}

func add(person Person) (int, error) {
	id := len(people) + 1
	people[id] = person
	return id, nil
}

func getByID(id int) (*Person, error) {
	person := people[id]
	if person.isEmpty() {
		return nil, errors.New("no person located by id")
	}
	return &person, nil
}
