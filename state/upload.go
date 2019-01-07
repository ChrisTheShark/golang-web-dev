package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFiles("templates/upload.gohtml"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var s string
		if req.Method == http.MethodPost {
			f, h, err := req.FormFile("file")
			if err != nil {
				log.Println(err)
				http.Error(w, "Unable to process file.", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			log.Printf("\n file: %v\nheader: %v\nerror: %v", f, h, err)

			bs, err := ioutil.ReadAll(f)
			if err != nil {
				log.Println(err)
				http.Error(w, "Unable to process file.", http.StatusInternalServerError)
				return
			}

			s = string(bs)
		}

		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, `
			<form method="post" enctype="multipart/form-data">
				<fieldset>
					<legend>Choose a File</legend>
					<label for="first_name">File to Upload:</label>
					<input type="file" id="file" name="file"><br>
					<input type="submit">
				</fieldset>
			</form>
		`+s)
	})
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
