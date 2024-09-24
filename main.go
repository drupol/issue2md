package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// GitHub API response structure for an issue
type Issue struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Number   int    `json:"number"`
	URL      string `json:"url"`
	Comments int    `json:"comments"`
	User     User   `json:"user"`
}

// GitHub API response structure for comments
type Comment struct {
	Body string `json:"body"`
	User User   `json:"user"`
}

type User struct {
	Login string `json:"login"`
}

// Parse GitHub Issue URL to get owner, repo, and issue number
func parseIssueURL(issueURL string) (owner, repo string, issueNumber int, err error) {
	parsedURL, err := url.Parse(issueURL)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid URL: %w", err)
	}

	// The URL path is in the format /{owner}/{repo}/issues/{issueNumber}
	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) != 4 || parts[2] != "issues" {
		return "", "", 0, fmt.Errorf("invalid GitHub issue URL format")
	}

	owner = parts[0]
	repo = parts[1]
	issueNumber, err = strconv.Atoi(parts[3])
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid issue number: %w", err)
	}

	return owner, repo, issueNumber, nil
}

// Fetch issue details using GitHub API
func fetchIssue(owner, repo string, issueNumber int, token string) (*Issue, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", owner, repo, issueNumber)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

// Fetch issue comments using GitHub API
func fetchComments(owner, repo string, issueNumber int, token string) ([]Comment, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/comments", owner, repo, issueNumber)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var comments []Comment
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, err
	}

	return comments, nil
}

// Convert issue and comments to markdown format
func issueToMarkdown(issue *Issue, comments []Comment) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n", issue.Title))
	sb.WriteString(fmt.Sprintf("**Issue Number**: #%d\n", issue.Number))
	sb.WriteString(fmt.Sprintf("**URL**: %s\n", issue.URL))
	sb.WriteString(fmt.Sprintf("**Created by**: %s\n\n", issue.User.Login))
	sb.WriteString(fmt.Sprintf("## Description\n%s\n\n", issue.Body))

	// Add comments section if there are any
	if len(comments) > 0 {
		sb.WriteString("## Comments\n")
		for i, comment := range comments {
			sb.WriteString(fmt.Sprintf("### Comment %d by %s\n", i+1, comment.User.Login))
			sb.WriteString(fmt.Sprintf("%s\n\n", comment.Body))
		}
	}

	return sb.String()
}

func main() {
	// Parse command-line arguments
	issueURL := flag.String("issue-url", "", "The GitHub issue URL")
	flag.Parse()

	if *issueURL == "" {
		fmt.Println("Error: issue URL is required")
		flag.Usage()
		return
	}

	// Parse the issue URL to get owner, repo, and issue number
	owner, repo, issueNumber, err := parseIssueURL(*issueURL)
	if err != nil {
		fmt.Printf("Error parsing issue URL: %v\n", err)
		return
	}

	token := os.Getenv("GITHUB_TOKEN") // GitHub token from environment variable

	// Fetch issue details
	issue, err := fetchIssue(owner, repo, issueNumber, token)
	if err != nil {
		fmt.Printf("Error fetching issue: %v\n", err)
		return
	}

	// Fetch issue comments
	comments, err := fetchComments(owner, repo, issueNumber, token)
	if err != nil {
		fmt.Printf("Error fetching comments: %v\n", err)
		return
	}

	// Convert issue and comments to markdown
	markdown := issueToMarkdown(issue, comments)

	// Save to file
	fileName := fmt.Sprintf("issue_%d.md", issue.Number)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = io.WriteString(file, markdown)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Issue and comments saved as Markdown in file %s\n", fileName)
}
