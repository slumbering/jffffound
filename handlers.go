package main

import (
	"net/http"
	"os"
	"path/filepath"
)

// homeHandler redirects root requests to the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Read the home page content from the static file
	content, err := os.ReadFile(filepath.Join(staticDir, "README.html"))
	if err != nil {
		http.Error(w, "Home page not found", http.StatusInternalServerError)
		return
	}

	// Set content type and serve the HTML content
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	pageName := r.URL.Path[len("/page/"):]
	if pageName == "" {
		http.NotFound(w, r)
		return
	}
	pagePath := staticDir + "/" + pageName
	content, err := os.ReadFile(pagePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
