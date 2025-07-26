package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v71/github"
)

const (
	repoOwner   = "slumbering"
	repoName    = "til"
	pageTimeout = 30 * time.Second
)

func newGithubClient() *github.Client {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		fmt.Println("GITHUB_TOKEN environment variable not set")
		return nil
	}

	return github.NewClient(nil).WithAuthToken(token)
}

func fetchRepoContent(client *github.Client, path string) (*Page, error) {
	ctx := context.Background()
	_, repoContent, resp, err := client.Repositories.GetContents(ctx, repoOwner, repoName, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching directory contents: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching directory contents: %s", resp.Status)
	}

	for _, content := range repoContent {
		if content.GetType() == "file" {
			file, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, content.GetPath(), nil)
			fileContent, err := file.GetContent()

			if err != nil {
				fmt.Println("Error extracting file content:", err)
				return nil, err
			}
			html, err := renderMarkdown(client, fileContent)

			if err != nil {
				fmt.Println("Error rendering markdown:", err)
				return nil, err
			}
			renderedHTML, err := renderWithLayout(*content.Name, html)
			if err != nil {
				return nil, fmt.Errorf("template rendering failed: %w", err)
			}
			newPage(*content.Name, renderedHTML)
		}
		if content.GetType() == "dir" {
			subPage, err := fetchRepoContent(client, content.GetPath())
			if err != nil {
				return nil, fmt.Errorf("error fetching subdirectory contents: %w", err)
			}
			if subPage != nil {
				fmt.Println("Subpage found:", subPage.Title)
				return subPage, nil
			}
		} else {
			fmt.Println("Skipping unsupported content type:", content.GetType(), "for", content.GetName())
		}
	}

	return nil, nil
}

// renderMarkdown converts markdown content to HTML
func renderMarkdown(client *github.Client, content string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pageTimeout)
	defer cancel()

	opt := &github.MarkdownOptions{
		Mode:    "gfm",
		Context: fmt.Sprintf("%s/%s", repoOwner, repoName),
	}

	output, _, err := client.Markdown.Render(ctx, content, opt)
	if err != nil {
		return "", fmt.Errorf("error rendering markdown: %w", err)
	}

	return output, nil
}

func getPages(client *github.Client) (*Page, error) {

	repoContent, err := fetchRepoContent(client, "")
	if err != nil {
		return nil, err
	}

	return repoContent, nil
}
