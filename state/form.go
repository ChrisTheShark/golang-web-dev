package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		value := req.FormValue("q")
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, `
			<form method="post">
				<input type="text" name="q">
				<input type="submit">
			</form>
		<br>`+value)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
