package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

type trival struct {
}

func (t *trival) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("Executing trival middleware")
	next(w, r)
}

func main() {
	r := mux.NewRouter()
	n := negroni.Classic()
	n.UseHandler(r)
	n.Use(&trival{})
	http.ListenAndServe(":8080", n)

}
