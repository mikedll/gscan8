
package main

import (
	"os"
	"net/http"
	"fmt"
	"log"
	"github.com/qor/render"
	"html/template"
	"encoding/json"
	"flag"
)

type TempParams struct {
	Name string
}

var sBootstrap template.HTML

var isProduction bool

func main() {
	isProduction = false

	flag.Parse()
	
	var addr string = ":8081"
	port := os.Getenv("PORT")
	appEnv := os.Getenv("APP_ENV")

	// Going to use this to determine production environment...LOL!
	if appEnv == "production" || port != "" {
		isProduction = true
		addr = fmt.Sprintf(":%s", port)
	}

	if flag.NArg() > 0 && flag.Arg(0) == "schema" {
		err := makeSchema(isProduction)
		if err != nil {
			log.Println("failed to create schema.")
			return
		}
		
		fmt.Println("created schema.")
		return
	}
	
	if flag.NArg() > 0 && flag.Arg(0) == "sample" {
		err := makeGists(isProduction)
		if err != nil {
			log.Println("failed to make sample gists", err)
			return
		}
		log.Println("sample gists created.")
		return
	}
	
	gists := getGistFiles()
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

	search := func(w http.ResponseWriter, req *http.Request) {
		// search db for json
		snippets := searchGistFiless(req.query)
		snippetsJson, err := json.Marshal(snippets)
		if err != nil {
			log.Println("error while marshalling snippets: ", err)
			snippetsJson = []byte{}
			// StatusInternalServerError           
			// write "Error while marshalling snippets.
			return
		}
		
		w.Header.Add("Content-Type", "application/json")
		w.write(string(snippetsJson))
		return
	}
	
	fmt.Println("Starting server...")
	http.Handle("/", http.HandlerFunc(root))

	http.Handle("/api/gists/search", http.HandlerFunc(search))

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
