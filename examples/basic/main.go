package main

import (
	"embed"
	"github.com/candlerb/tiny_jsonrpc"
	"html/template"
	"log"
	"net/http"
)

//go:embed html/*.html static/*.js
var content embed.FS
var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFS(content, "html/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: ", err)
	}
}

func main() {
	http.HandleFunc("/", PageIndex)
	http.HandleFunc("/rpc", MyHandler.HTTPHandler)
	http.Handle("/static/", http.FileServer(http.FS(content)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func PageIndex(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.html", struct{}{})
	if err != nil {
		log.Print("ExecuteTemplate error: ", err)
		rpc.HTTPError(w, http.StatusInternalServerError, err)
	}
}
