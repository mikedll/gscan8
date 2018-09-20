
package main

import (
	"flag"
	"net/http"
	"fmt"
	"log"
	"html/template"
	"github.com/mikedll/stringutil"
)

var addr = flag.String("addr", ":8081", "http service address")

var templ = template.Must(template.New("someserver").Parse("Hey mike. We got some fanciness in here: 語語語 {{.}}."))

func main() {
	fmt.Println(stringutil.Reverse("!oG olleH"))
	fmt.Println("Unhex:", unhex('a'))
	fmt.Println("Unhex:", unhex('F'))
	fmt.Println("Unhex:", unhex('c'))
	fmt.Println("Unhex:", unhex('C'))
	fmt.Println("Unhex:", unhex('8'))

	for _, char := range "日本\x80語ab" {
		fmt.Printf("here is a char, %#U, starting.\n", char)
	}

	http.Handle("/", http.HandlerFunc(root))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func root(w http.ResponseWriter, req *http.Request) {
	templ.Execute(w, "Carol")
}

func unhex(b byte) byte {
	switch {
	case '0' <= b && b <= '9':
		return b - '0'
	case 'a' <= b && b <= 'f':
		return b - 'a' + 10
	case 'A' <= b && b <= 'F':
		return b - 'A' + 10
	}

	return 0;
}
