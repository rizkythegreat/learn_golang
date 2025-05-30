package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

type Info struct {
    Affiliation string
    Address string
}

type Person struct {
    Name string
    Gender string
    Hobbies []string
    Info Info
}

func (t Info) GetAffiliationDetailInfo() string {
    return "have 31 divisions"
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        var person = Person{
            Name:    "Bruce Wayne",
            Gender:  "male",
            Hobbies: []string{"Reading Books", "Traveling", "Buying things"},
            Info:    Info{"Wayne Enterprises", "Gotham City"},
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