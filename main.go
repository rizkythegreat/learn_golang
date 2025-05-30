package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        var filePath = path.Join("views", "index.html")
        var tmpl, err = template.ParseFiles(filePath)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        var data = map[string]interface{}{
            "title": "Learning Golang Web",
            "name": "Batman",
        }

        err = tmpl.Execute(w, data)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    http.Handle("/static/", 
        http.StripPrefix("/static/", 
            http.FileServer(http.Dir("assets"))))
    fmt.Println("server started at localhost:9000")
    http.ListenAndServe(":9000", nil)
}