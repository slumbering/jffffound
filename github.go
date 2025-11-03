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

type Document struct {
	isCategory bool
	title      string
	path       string
	content    *string
}

func scanRepo(client *github.Client, path string) ([]*github.RepositoryContent, error) {
	ctx := context.Background()
	_, repoContent, resp, err := client.Repositories.GetContents(ctx, repoOwner, repoName, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching directory contents: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching directory contents: %s", resp.Status)
	}

	return repoContent, nil
}

func prepareDocuments(client *github.Client, repoContent []*github.RepositoryContent) ([]Document, error) {
	ctx := context.Background()
	document := make([]Document, 0)
	for _, content := range repoContent {
		if content.GetType() == "dir" {
			document = append(document, Document{
				isCategory: true,
				title:      content.GetName(),
				content:    nil,
			})

			// Recursively scan subdirectory
			subRepoContent, err := scanRepo(client, content.GetPath())
			if err != nil {
				return nil, fmt.Errorf("error scanning subdirectory: %w", err)
			}

			// Get subdirectory documents and append them
			subDocuments, err := prepareDocuments(client, subRepoContent)
			if err != nil {
				return nil, err
			}
			document = append(document, subDocuments...)
		}

		if content.GetType() == "file" {
			file, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, content.GetPath(), nil)
			if err != nil {
				return nil, fmt.Errorf("error fetching file content: %w", err)
			}
			fileContent, err := file.GetContent()
			if err != nil {
				return nil, fmt.Errorf("error extracting file content: %w", err)
			}
			html, err := renderMarkdown(client, fileContent)
			if err != nil {
				return nil, fmt.Errorf("error rendering markdown: %w", err)
			}

			document = append(document, Document{
				isCategory: false,
				title:      content.GetName(),
				path:       content.GetPath(),
				content:    &html,
			})
		}
	}

	return document, nil
}

func renderPages(client *github.Client, path string) (*Document, error) {
	repoScanned, err := scanRepo(client, path)
	if err != nil {
		return nil, fmt.Errorf("error scanning repository: %w", err)
	}

	preparedDocuments, err := prepareDocuments(client, repoScanned)
	fmt.Println("Prepared Documents:", preparedDocuments)
	if err != nil {
		return nil, fmt.Errorf("error preparing pages: %w", err)
	}

	renderWithLayout(preparedDocuments)

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

func getPages(client *github.Client) (*Document, error) {
	repoContent, err := renderPages(client, "")
	if err != nil {
		return nil, err
	}

	return repoContent, nil
}
