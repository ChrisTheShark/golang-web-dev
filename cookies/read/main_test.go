package main

import (
	"io/ioutil"
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
	assert.Equal(t, "Session_ID", cookies[0].Name)
	assert.Equal(t, "1234567890", cookies[0].Value)
}

func TestGet(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/get", nil)
	w := httptest.NewRecorder()
	r.AddCookie(&http.Cookie{
		Name:  "Session_ID",
		Value: "1234567890",
	})

	get(w, r)
	resp := w.Result()
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "Session_ID=1234567890", string(bs))
}
