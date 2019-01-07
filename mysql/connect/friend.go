package main

import (
	"database/sql"
	"log"
)

func getByID(id int) (*Friend, error) {
	var name string
	row := db.QueryRow("SELECT name FROM friends WHERE friend_id = ?", id)
	if err := row.Scan(&name); err != nil {
		return nil, err
	}
	return &Friend{
		ID:   id,
		Name: name,
	}, nil
}

func deleteByID(id int) error {
	rs, err := db.Exec("DELETE FROM friends WHERE friend_id = ?", id)
	if err != nil {
		return err
	}

	num, err := rs.RowsAffected()
	if err != nil {
		return err
	} else if num != 1 {
		return sql.ErrNoRows
	}
}

func getAll() ([]Friend, error) {
	rows, err := db.Query("SELECT friend_id, name FROM friends")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var friends []Friend
	for rows.Next() {
		var friend Friend
		if err := rows.Scan(&friend.ID, &friend.Name); err != nil {
			log.Printf("error scanning row, continuing: %v", err)
			continue
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func add(name string) (*Friend, error) {
	rs, err := db.Exec("INSERT INTO friends (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Friend{
		ID:   int(id),
		Name: name,
	}, nil
}
