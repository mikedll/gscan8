
package main

import (
	"os"
	"net/http"
	"fmt"
	"log"
	"github.com/qor/render"
	"html/template"
	"encoding/json"
)

type TempParams struct {
	Name string
}

var sBootstrap template.HTML

func main() {
	
	var addr string = ":8081"
	port := os.Getenv("PORT")

	if port != "" {
		addr = fmt.Sprintf(":%s", port)
	}
	
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
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
