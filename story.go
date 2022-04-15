package choosebook

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

var defaultHandlerTempl = `
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
        <li><a href="/{{.Chapter}}">{{.Text}}</a> </li>    
        {{end}}
        </ul>
    </body>
</html>`

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type MyHandler struct {
	story    Story
	template *template.Template
	pathFunc func(r *http.Request) string
}

type HandlerOption func(r *MyHandler)

type Story map[string]Chapter

func JsonStory(r io.Reader) (Story, error) {
	var story Story

	d := json.NewDecoder(r)

	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *MyHandler) {
		h.template = t
	}
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *MyHandler) {
		h.pathFunc = fn
	}
}

func NewHandler(s Story, opts ...HandlerOption) MyHandler {
	h := MyHandler{s, tpl, defaultPathFunc}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

func defaultPathFunc(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:]
}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFunc(r)

	if chapter, ok := h.story[path]; ok {
		err := h.template.Execute(w, chapter)

		if err != nil {
			log.Printf("%v", err)

			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found...", http.StatusNotFound)
}

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTempl))
}
