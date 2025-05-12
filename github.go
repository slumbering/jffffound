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

// newGithubClient creates a new GitHub client with token from environment
func newGithubClient() *github.Client {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		fmt.Println("GITHUB_TOKEN environment variable not set")
		return nil
	}

	return github.NewClient(nil).WithAuthToken(token)
}

func fetchReadmeContent(client *github.Client) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pageTimeout)
	defer cancel()

	readmeContent, _, err := client.Repositories.GetReadme(ctx, repoOwner, repoName, nil)
	if err != nil {
		return "", fmt.Errorf("error fetching repository contents: %w", err)
	}

	content, err := readmeContent.GetContent()
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

// TODO: This function should fetch every pages from the root and subdirectories
// If the file is equal to README it should be considered as the home page
// Otherwise, it should be considered as a subpage.
func getPages(client *github.Client) (*Page, error) {

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
