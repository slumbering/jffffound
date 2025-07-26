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

	_, err := getPages(githubClient)
	if err != nil {
		log.Fatalf("Failed to generate homepage: %v", err)
	}

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/page/", pageHandler)

	// Start the server
	log.Printf("Server listening on port %s", serverPort)
	err = http.ListenAndServe(serverPort, nil)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
