
package main

import (
	"flag"
	"net/http"
	"fmt"
	"log"
	"html/template"
)



var addr = flag.String("addr", ":8081", "http service address")

type TempParams struct {
	Name string
}

func main() {
	fmt.Println("Starting server...")

	var templ, _ = template.ParseFiles("index.html")
	var root = func(w http.ResponseWriter, req *http.Request) {
		m := TempParams{Name: "Carol"}
		templ.Execute(w, m)
	}

	http.Handle("/", http.HandlerFunc(root))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
