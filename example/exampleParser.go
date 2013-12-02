package main

import (
	"fmt"
	"github.com/airdispatch/dpl"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	xmlFile, err := os.Open("example/example.dpl")
	if err != nil {
		fmt.Println("Couldn't open file.")
		return
	}
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)

	o, err := dpl.ParseDPL(b)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Hello, %q", r.URL.Path)
		action := r.URL.Path[1:]
		a := o.Actions[action]
		fmt.Fprintf(w, "%s", a.HTML)
	})
	log.Fatal(http.ListenAndServe(":2048", nil))
}
