package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var tmpl = template.Must(template.ParseFiles("base.html"))

func rootHandler(w http.ResponseWriter, r *http.Request) {
	dat := struct {
		Title string
		Time  time.Time
	}{
		Title: "Test",
		Time:  time.Now(),
	}
	if err := tmpl.ExecuteTemplate(w, "base.html", dat); err != nil {
		log.Fatal(err)
	}
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
