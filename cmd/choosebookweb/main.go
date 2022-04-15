package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/evgeniy-dammer/choosebook"
)

var storyTempl = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Choose Your Own Adventure Book</title>
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
        <p>{{.}}</p>    
        {{end}}
        <ul>
        {{range .Options}}
        <li><a href="/story/{{.Chapter}}">{{.Text}}</a> </li>    
        {{end}}
        </ul>
    </body>
</html>`

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)

	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}

	return path[len("/story/"):]
}

func main() {
	port := flag.Int("port", 3000, "port to start web application on")
	file := flag.String("file", "gopher.json", "JSON file with chapters")
	flag.Parse()

	fmt.Printf("Using chapters from %s\n", *file)

	f, err := os.Open(*file)

	if err != nil {
		panic(err)
	}

	story, err := choosebook.JsonStory(f)

	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.New("").Parse(storyTempl))

	h := choosebook.NewHandler(story, choosebook.WithTemplate(tmpl), choosebook.WithPathFunc(pathFn))

	mux := http.NewServeMux()
	mux.Handle("/story/", h)

	fmt.Printf("Starting Web Application at %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}
