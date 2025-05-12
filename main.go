// main.go
package main

import (
	"log"
	"net/http"
)

const (
	staticDir  = "./static"
	serverPort = ":3000"
)

func main() {
	// Initialize GitHub client
	githubClient := newGithubClient()

	page, err := getPages(githubClient)
	if err != nil {
		log.Fatalf("Failed to generate homepage: %v", err)
	}
	log.Printf("Generated %s.html successfully", page.Title)

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)

	// Start the server
	log.Printf("Server listening on port %s", serverPort)
	err = http.ListenAndServe(serverPort, nil)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
