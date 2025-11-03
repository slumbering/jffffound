// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

const (
	staticDir  = "./static"
	serverPort = ":3000"
)

func main() {
	// Ensure static directory exists
	createStaticDir()
	// Copy CSS directory to static
	copyCSSDir()

	// Initialize GitHub client
	githubClient := newGithubClient()

	_, err := getPages(githubClient)
	if err != nil {
		log.Fatalf("Failed to generate homepage: %v", err)
	}

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", pageHandler)

	// Start the server
	log.Printf("Server listening on port %s", serverPort)
	err = http.ListenAndServe(serverPort, nil)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

/**
 * createStaticDir ensures the static directory exists
 */
func createStaticDir() error {
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		return fmt.Errorf("failed to create static directory: %w", err)
	}
	return nil
}

func copyCSSDir() error {
	sourceDir := "./css"
	destDir := filepath.Join(staticDir, "css")
	err := cp.Copy(sourceDir, destDir)
	if err != nil {
		return fmt.Errorf("failed to copy CSS directory: %w", err)
	}

	return nil
}
