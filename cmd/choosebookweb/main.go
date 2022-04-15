package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/evgeniy-dammer/choosebook"
)

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

	h := choosebook.NewHandler(story)
	fmt.Printf("Starting Web Application at %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
