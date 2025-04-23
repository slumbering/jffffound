package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v71/github"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".html"
	return os.WriteFile(filepath.Join("./static", filename), p.Body, 0600)
}

func getHomePage() *Page {
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))

	fileContent, _, _, error := client.Repositories.GetContents(context.Background(), "slumbering", "til", "README.md", nil)
	if error != nil {
		fmt.Println("Error fetching organizations:", error)
		return nil
	}
	content, err := fileContent.GetContent()
	if err != nil {
		fmt.Println("Error getting content:", err)
		return nil
	}
	opt := &github.MarkdownOptions{
		Mode:    "gfm",
		Context: "slumbering/til",
	}

	output, _, err := client.Markdown.Render(context.Background(), content, opt)
	if err != nil {
		fmt.Println("Error rendering markdown:", err)
		return nil
	}
	p := &Page{
		Title: "home",
		Body:  []byte(output),
	}

	p.save()

	return p
}

// func viewHandler(w http.ResponseWriter, r *http.Request) {
// 	source := []byte("# Hello World\nThis is a sample markdown content.")
// 	var buf bytes.Buffer
// 	if err := goldmark.Convert(source, &buf); err != nil {
// 		panic(err)
// 	}
// 	p := &Page{
// 		Body: buf.Bytes(),
// 	}
// 	w.Header().Set("Content-Type", "text/html")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(string(p.Body)))
// }

func main() {
	getHomePage()
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		log.Fatal(err)
	}

	// http.HandleFunc("/", getHomePage)
	// http.HandleFunc("/til/", viewHandler)
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	fmt.Println("Error starting server:", err)
	// }
}
