package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

func routeIndexGet(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        var filePath = path.Join("views", "index.html")
        var tmpl, errTemplate = template.New("form").ParseFiles(filePath)
        if errTemplate != nil {
            http.Error(w, errTemplate.Error(), http.StatusInternalServerError)
            return
        }
        var err = tmpl.Execute(w, nil)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    http.Error(w, "", http.StatusBadRequest)
}
func routeSubmitPost(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        var filePath = path.Join("views", "index.html")
        var tmpl = template.Must(template.New("result").ParseFiles(filePath))

        if err := r.ParseForm(); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        var name = r.FormValue("name")
        var message = r.Form.Get("message")

        var data = map[string]string{"name": name, "message": message}

        if err := tmpl.Execute(w, data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    http.Error(w, "", http.StatusBadRequest)
}

func main() {
    http.HandleFunc("/", routeIndexGet)
    http.HandleFunc("/process", routeSubmitPost)
    fmt.Println("server started at localhost:9000")
    http.ListenAndServe(":9000", nil)
}