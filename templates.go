package main

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Title   string
	Content template.HTML
	Pages   map[string]string // Map of page names to file names
}

func renderWithLayout(title string, content string) ([]byte, error) {
	// Parse the manually created layout template
	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		return nil, err
	}

	//list of html pages within the static directory
	pagePaths := staticDir + "/" // List all HTML files in the static directory
	files, err := os.ReadDir(pagePaths)
	if err != nil {
		return nil, err
	}
	// Create a map to hold the page names
	pageNames := make(map[string]string)
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".html" && file.Name() != "README.html" {
			// Use the file name without extension as the key
			pageName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			pageNames[pageName] = file.Name() // Store the full file name
		}
	}
	// Prepare template data
	data := TemplateData{
		Title:   title,
		Content: template.HTML(content),
		Pages:   pageNames,
	}

	// Render the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
