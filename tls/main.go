package main

import (
	"log"
	"net/http"
)

// Generate Self-Signed Certificate:
// go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=localhost

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text-plain")
		w.Write([]byte("This is an example server."))
	})
	log.Fatal(http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil))
}
