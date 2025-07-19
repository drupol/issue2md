package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/bigwhite/issue2md/internal/converter"
	"github.com/bigwhite/issue2md/internal/github"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	issueURL := r.FormValue("issue_url")
	if issueURL == "" {
		http.Error(w, "Issue URL is required", http.StatusBadRequest)
		return
	}

	enableReactions := r.FormValue("enable_reactions") == "true"
	enableUserLinks := r.FormValue("enable_user_links") == "true"

	owner, repo, issueNumber, issueType, err := github.ParseURL(issueURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing issue URL: %v", err), http.StatusBadRequest)
		return
	}

	token := os.Getenv("GITHUB_TOKEN")

	var markdown string
	var filename string

	switch issueType {
	case "issue":
		issue, err := github.FetchIssue(owner, repo, issueNumber, token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching issue: %v", err), http.StatusInternalServerError)
			return
		}

		comments, err := github.FetchComments(owner, repo, issueNumber, token, enableReactions, enableUserLinks)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching comments: %v", err), http.StatusInternalServerError)
			return
		}
		markdown = converter.IssueToMarkdown(issue, comments, enableUserLinks)
		filename = fmt.Sprintf("%s_%s_issue_%d.md", owner, repo, issue.Number)
	case "discussion":
		discussion, err := github.FetchDiscussion(owner, repo, issueNumber, token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching discussion: %v", err), http.StatusInternalServerError)
			return
		}
		discussionComments, err := github.FetchDiscussionComments(owner, repo, issueNumber, token, enableReactions)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching discussion comments: %v", err), http.StatusInternalServerError)
			return
		}
		markdown = converter.DiscussionToMarkdown(discussion, discussionComments, enableUserLinks)
		filename = fmt.Sprintf("%s_%s_discussion_%d.md", owner, repo, discussion.Number)
	default:
		http.Error(w, fmt.Sprintf("Unsupported URL type: %s", issueType), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/markdown")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	_, err = w.Write([]byte(markdown))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
	}
}
