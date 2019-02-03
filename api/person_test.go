package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
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

	actual, err := getAll()

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, 1, len(actual))
}

func TestGetAllEmpty(t *testing.T) {
	people = map[int]Person{}

	actual, err := getAll()

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, 0, len(actual))
}

func TestGeyByID(t *testing.T) {
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

	person, err := getByID(1)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, people[1], *person)
}

func TestGeyByIDNotFound(t *testing.T) {
	people = map[int]Person{}

	person, err := getByID(1)

	assert.NotNil(t, err, "error should not be nil")
	assert.Equal(t, "no person located by id", err.Error())
	assert.Nil(t, person, "error should be nil")
}

func TestAdd(t *testing.T) {
	people = map[int]Person{}

	person1 := Person{
		Name: "Amy Jonson",
		Age:  39,
		Friends: []string{
			"Jak",
			"Daxter",
		},
	}

	person2 := Person{
		Name: "Amy Jonson",
		Age:  39,
		Friends: []string{
			"Jak",
			"Daxter",
		},
	}

	id1, err := add(person1)
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, 1, id1)

	id2, err := add(person2)
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, 2, id2)

	assert.Equal(t, 2, len(people))
}
