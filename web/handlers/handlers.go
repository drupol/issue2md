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

	owner, repo, issueNumber, err := github.ParseIssueURL(issueURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing issue URL: %v", err), http.StatusBadRequest)
		return
	}

	token := os.Getenv("GITHUB_TOKEN")

	issue, err := github.FetchIssue(owner, repo, issueNumber, token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching issue: %v", err), http.StatusInternalServerError)
		return
	}

	comments, err := github.FetchComments(owner, repo, issueNumber, token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching comments: %v", err), http.StatusInternalServerError)
		return
	}

	markdown := converter.IssueToMarkdown(issue, comments)

	w.Header().Set("Content-Type", "text/markdown")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%s_issue_%d.md", owner, repo, issue.Number))
	_, err = w.Write([]byte(markdown))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
	}
}
