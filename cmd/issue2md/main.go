package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bigwhite/issue2md/internal/converter"
	"github.com/bigwhite/issue2md/internal/github"
	flag "github.com/spf13/pflag"
)

// This will be set at build time
var version = "dev"

var (
	enableReactions = flag.Bool("enable-reactions", false, "Include reactions in the output.")
	showVersion     = flag.Bool("version", false, "Print the version and exit")
)

func usage() {
	fmt.Println("Usage: issue2md [--enable-reactions] [--version] issue-url [markdown-file]")
	fmt.Println("Arguments:")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	// Check for the version flag first.
	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		fmt.Println("Error: issue-url is required.")
		usage()
		return
	}

	issueURL := flag.Arg(0)
	var markdownFile string
	if len(flag.Args()) >= 2 {
		markdownFile = flag.Arg(1)
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

		comments, err := github.FetchComments(owner, repo, issueNumber, token, *enableReactions)
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

		discussionComments, err := github.FetchDiscussionComments(owner, repo, issueNumber, token, *enableReactions)
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
