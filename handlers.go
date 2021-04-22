package main

import (
	"net/http"
)

func index(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h4>Server up and running</h4>"))
}
