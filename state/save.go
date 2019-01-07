package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFiles("templates/upload.gohtml"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			f, h, err := req.FormFile("file")
			if err != nil {
				log.Println(err)
				http.Error(w, "Unable to locate file.", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			log.Printf("\n file: %v\nheader: %v\nerror: %v", f, h, err)

			bs, err := ioutil.ReadAll(f)
			if err != nil {
				log.Println(err)
				http.Error(w, "Unable to read file.", http.StatusInternalServerError)
				return
			}

			dst, err := os.Create(filepath.Join("./user", h.Filename))
			if err != nil {
				log.Println(err)
				http.Error(w, "Unable to save file.", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			_, err = dst.Write(bs)
			if err != nil {
				log.Println(err)
				http.Error(w, "Unable to write contents to file.", http.StatusInternalServerError)
				return
			}
		}

		templates.Execute(w, nil)
	})
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
