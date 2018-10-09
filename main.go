
package main

import (
	"flag"
	"net/http"
	"fmt"
	"log"
	"github.com/qor/render"
	"html/template"
)

var addr = flag.String("addr", ":8081", "http service address")

type TempParams struct {
	Name string
}

var sBootstrap template.HTML

func main() {
	main2()
	
	sBootstrap = template.HTML(`[
        {title: "A neat gist", href: "https://gist.github.com/mikedll/8eaa6df25ac7a10ae3ded33e7f00b306"},
        {title: "Cute Gist", href: "https://gist.github.com/mikedll/db0bbe17ddfa389eada54682f4a5b4c5"}
      ]`)

	fmt.Println("Starting server...")

	Render := render.New(&render.Config{
		ViewPaths:  []string{},
		DefaultLayout: "",
		FuncMapMaker: nil,
	})

	root := func(w http.ResponseWriter, req *http.Request) {
		ctx := map[string]template.HTML{"Bootstrap": sBootstrap}
		Render.Execute("index", ctx, req, w)
	}
	
	http.Handle("/", http.HandlerFunc(root))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
