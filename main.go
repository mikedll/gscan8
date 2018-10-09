
package main

import (
	"flag"
	"net/http"
	"fmt"
	"log"
	"github.com/qor/render"
	"html/template"
	"encoding/json"
)

var addr = flag.String("addr", ":8081", "http service address")

type TempParams struct {
	Name string
}

var sBootstrap template.HTML

func main() {
	makeGists()

	gists := getGists()
	gistsJson, err := json.Marshal(gists)
	if err != nil {
		log.Println("unable to find gists: ", err)
		gistsJson = []byte{}
	}

	Render := render.New(&render.Config{
		ViewPaths:  []string{},
		DefaultLayout: "",
		FuncMapMaker: nil,
	})

	root := func(w http.ResponseWriter, req *http.Request) {
		sBootstrap = template.HTML(string(gistsJson))
		ctx := map[string]template.HTML{"Bootstrap": sBootstrap}
		Render.Execute("index", ctx, req, w)
	}
	
	fmt.Println("Starting server...")
	http.Handle("/", http.HandlerFunc(root))
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
