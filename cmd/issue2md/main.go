package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bigwhite/issue2md/internal/converter"
	"github.com/bigwhite/issue2md/internal/github"
)

func usage() {
	fmt.Println("Usage: issue2md [--enable-reactions] issue-url [markdown-file]")
	fmt.Println("Arguments:")
	fmt.Println("  --enable-reactions  (optional) Include reactions in the output.")
	fmt.Println("  issue-url           The URL of the GitHub issue to convert.")
	fmt.Println("  markdown-file       (optional) The output markdown file.")
}

func main() {
	enableReactions := false
	args := os.Args[1:]
	// Check for --enable-reactions flag
	for i, arg := range args {
		if arg == "--enable-reactions" {
			enableReactions = true
			// Remove the flag from args
			args = append(args[:i], args[i+1:]...)
			break
		}
	}

	if len(args) < 1 {
		fmt.Println("Error: issue-url is required.")
		usage()
		return
	}

	issueURL := args[0]
	var markdownFile string
	if len(args) >= 2 {
		markdownFile = args[1]
	}

	owner, repo, issueNumber, issueType, err := github.ParseURL(issueURL)
	if err != nil {
		fmt.Printf("Error parsing issue URL: %v\n", err)
		return
	}

	token := os.Getenv("GITHUB_TOKEN")

	var markdown string
	switch issueType {
	case "issue":
		issue, err := github.FetchIssue(owner, repo, issueNumber, token)
		if err != nil {
			fmt.Printf("Error fetching issue: %v\n", err)
			return
		}

		comments, err := github.FetchComments(owner, repo, issueNumber, token, enableReactions)
		if err != nil {
			fmt.Printf("Error fetching comments: %v\n", err)
			return
		}
		markdown = converter.IssueToMarkdown(issue, comments)

	case "discussion":
		discussion, err := github.FetchDiscussion(owner, repo, issueNumber, token)
		if err != nil {
			fmt.Printf("Error fetching discussion: %v\n", err)
			return
		}

		discussionComments, err := github.FetchDiscussionComments(owner, repo, issueNumber, token, enableReactions)
		if err != nil {
			fmt.Printf("Error fetching discussion comments: %v\n", err)
			return
		}
		markdown = converter.DiscussionToMarkdown(discussion, discussionComments)

	default:
		fmt.Printf("Unsupported URL type: %s\n", issueType)
		return
	}

	if markdownFile == "" {
		if _, err := io.WriteString(os.Stdout, markdown); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to stdout: %v\n", err)
			os.Exit(1)
		}
		return
	}

	file, err := os.Create(markdownFile)
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

	fmt.Printf("Issue and comments saved as Markdown in file %s\n", markdownFile)
}
