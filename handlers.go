package main

import (
	"net/http"
	"text/template"
	"time"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("./templates/index.html"))
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tpl.Execute(w, time.Now().Year())
}
