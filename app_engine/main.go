package main

import (
	"net/http"
	"google.golang.org/appengine"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	appengine.Main()
}
