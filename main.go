
package main

import (
	"flag"
	"net/http"
	"fmt"
	"log"
	"github.com/qor/render"
)



var addr = flag.String("addr", ":8081", "http service address")

type TempParams struct {
	Name string
}

func main() {
	fmt.Println("Starting server...")

	Render := render.New(&render.Config{
		ViewPaths:  []string{},
		DefaultLayout: "",
		FuncMapMaker: nil,
	})	
	
	var root = func(w http.ResponseWriter, req *http.Request) {
		Render.Execute("index", nil, req, w)
	}

	http.Handle("/", http.HandlerFunc(root))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
