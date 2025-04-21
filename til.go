//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/yuin/goldmark"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	source := []byte("# Hello World\nThis is a sample markdown content.")
	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		panic(err)
	}
	p := &Page{
		Body: buf.Bytes(),
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(p.Body)))
}

func main() {
	http.HandleFunc("/til/", viewHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
