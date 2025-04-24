package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Page represents an HTML page with a title and body content
type Page struct {
	Title string
	Body  []byte
}

// save writes the page content to a file in the static directory
func (p *Page) save() error {
	filename := p.Title + ".html"
	filePath := filepath.Join(staticDir, filename)

	// Ensure the static directory exists
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		return fmt.Errorf("failed to create static directory: %w", err)
	}

	return os.WriteFile(filePath, p.Body, 0600)
}

// newPage creates and saves a new page
func newPage(title string, body []byte) (*Page, error) {
	page := &Page{
		Title: title,
		Body:  body,
	}

	if err := page.save(); err != nil {
		return nil, fmt.Errorf("failed to save page: %w", err)
	}

	return page, nil
}
