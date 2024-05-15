package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const defaultPort = "80"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	port, ok := os.LookupEnv("PORT")
	if !ok {
		fmt.Println("port not specified in environment variable.")
		fmt.Printf("using default port of %s. \n", port)
		port = defaultPort
	}

	fmt.Printf("starting front end service on port %s.", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {

	htmlTemplatePath, ok := os.LookupEnv("HTML_TEMPLATES_PATH")
	if !ok {
		htmlTemplatePath = "./cmd/web/templates/"
	}

	partials := []string{
		htmlTemplatePath + "/base.layout.gohtml",
		htmlTemplatePath + "/header.partial.gohtml",
		htmlTemplatePath + "/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf(htmlTemplatePath+"/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
