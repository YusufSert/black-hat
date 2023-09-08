package main

import (
    "fmt"
    "net/http"
)

type router struct {
}

func (r*router) ServeHTTP(w http.ResponseWriter, req*http.Request) {
    switch req.URL.Path {
    case "/a":
        fmt.Fprint(w, "Exucuting /a")
    case "/b":
        fmt.Fprint(w, "Exucuting /b")
    case "/c":
        fmt.Fprint(w, "Exucuting /c")
    default:
        http.Error(w, "404 Not Found", 404)
    }
}

func main() {
    var r router
    http.ListenAndServe(":8000", &r)
}
