package main

import (
	"bytes"
	"html/template"
	"strings"
)

type Menu struct {
	Category string
	Item     string
	Path     string
}

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Title   string
	Content template.HTML
	Menu    []Menu
}

func renderWithLayout(documents []Document) (*Page, error) {
	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		return nil, err
	}

	menu := buildMenu(documents)

	for _, document := range documents {
		if document.isCategory {
			continue
		}

		data := TemplateData{
			Title: document.title,
			Menu:  menu,
		}

		if document.content != nil {
			data.Content = template.HTML(*document.content)
		}

		// Render the template
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return nil, err
		}

		newPage(document.title, buf.Bytes())
	}

	return nil, nil
}

func buildMenu(documents []Document) []Menu {
	menu := make([]Menu, 0)
	var currentCategory string

	for _, doc := range documents {
		if doc.isCategory {
			// It's a category/directory
			currentCategory = doc.title
			menu = append(menu, Menu{
				Category: currentCategory,
				Item:     "", // No item for category header
			})
		} else {
			// It's a file/item
			itemPath := strings.TrimSuffix(doc.title, ".md") + ".html"
			itemName := strings.TrimSuffix(doc.title, ".md")
			menu = append(menu, Menu{
				Category: currentCategory, // Group under current category
				Item:     itemName,
				Path:     itemPath,
			})
		}
	}

	return menu
}
