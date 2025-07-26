package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func pageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
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

	pageName := r.URL.Path[len("/"):]
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
