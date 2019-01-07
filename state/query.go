package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		value := req.FormValue("q")
		io.WriteString(w, value)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
