package main

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func createSession() (*http.Cookie, error) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: sessionID.String(),
	}
	return cookie, nil
}

func hasActiveSession(r *http.Request) bool {
	_, err := r.Cookie("session")
	return err == nil
}
