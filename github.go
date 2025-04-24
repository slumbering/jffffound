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
	readmeFile  = "README.md"
	pageTimeout = 30 * time.Second
)

// newGitHubClient creates a new GitHub client with token from environment
func newGitHubClient() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	return github.NewClient(nil).WithAuthToken(token)
}

// fetchReadmeContent retrieves README.md content from GitHub
func fetchReadmeContent(client *github.Client) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pageTimeout)
	defer cancel()

	fileContent, _, _, err := client.Repositories.GetContents(
		ctx,
		repoOwner,
		repoName,
		readmeFile,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("error fetching repository contents: %w", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("error extracting content: %w", err)
	}

	return content, nil
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

func getHomePage(client *github.Client) (*Page, error) {
	// Fetch and render the markdown content
	content, err := fetchReadmeContent(client)
	if err != nil {
		return nil, err
	}

	html, err := renderMarkdown(client, content)
	if err != nil {
		return nil, err
	}

	// Wrap the HTML content in our layout template
	renderedHTML, err := renderWithLayout("Today I Learned", html)
	if err != nil {
		return nil, fmt.Errorf("template rendering failed: %w", err)
	}

	return newPage("home", renderedHTML)
}
