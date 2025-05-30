package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

type Superhero struct {
    Name    string
    Alias   string
    Friends []string
}

func (s Superhero) SayHello(from string, message string) string {
    return fmt.Sprintf("%s said: \"%s\"", from, message)
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        var person = Superhero{
            Name:    "Bruce Wayne",
            Alias:   "Batman",
            Friends: []string{"Superman", "Flash", "Green Lantern"},
        }

        var filePath = path.Join("views", "index.html")
        var tmpl, errTemplate = template.ParseFiles(filePath)
        if errTemplate != nil {
            http.Error(w, errTemplate.Error(), http.StatusInternalServerError)
            return
        }
        if err := tmpl.Execute(w, person); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    fmt.Println("server started at localhost:9000")
    http.ListenAndServe(":9000", nil)
}