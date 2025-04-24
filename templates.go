package main

import (
	"bytes"
	"html/template"
)

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Title   string
	Content template.HTML
}

func renderWithLayout(title string, content string) ([]byte, error) {
	// Parse the manually created layout template
	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		return nil, err
	}

	// Prepare template data
	data := TemplateData{
		Title:   title,
		Content: template.HTML(content),
	}

	// Render the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
