package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	index(w, r)
	resp := w.Result()

	cookies := resp.Cookies()

	assert.Equal(t, "counter", cookies[0].Name)
	assert.Equal(t, "1", cookies[0].Value)
}

func TestReset(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/reset", nil)
	r.AddCookie(&http.Cookie{
		Name:  "counter",
		Value: "1",
	})
	w := httptest.NewRecorder()

	reset(w, r)
	resp := w.Result()

	cookies := resp.Cookies()

	assert.Equal(t, "counter", cookies[0].Name)
	assert.Equal(t, "1", cookies[0].Value)
	assert.Equal(t, -1, cookies[0].MaxAge)
}
