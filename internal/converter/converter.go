package converter

import (
	"fmt"
	"strings"

	"github.com/bigwhite/issue2md/internal/github"
)

func IssueToMarkdown(issue *github.Issue, comments []github.Comment) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n", issue.Title))
	sb.WriteString(fmt.Sprintf("**Issue Number**: #%d\n", issue.Number))
	sb.WriteString(fmt.Sprintf("**URL**: %s\n", issue.URL))
	sb.WriteString(fmt.Sprintf("**Created by**: %s\n\n", issue.User.Login))
	sb.WriteString(fmt.Sprintf("## Description\n%s\n\n", issue.Body))

	if len(comments) > 0 {
		sb.WriteString("## Comments\n")
		for i, comment := range comments {
			sb.WriteString(fmt.Sprintf("### Comment %d by %s\n", i+1, comment.User.Login))
			sb.WriteString(fmt.Sprintf("%s\n\n", comment.Body))
		}
	}

	return sb.String()
}

func DiscussionToMarkdown(discussion *github.Discussion, discussionComments []github.DiscussionComment) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n", discussion.Title))
	sb.WriteString(fmt.Sprintf("**Discussion Number**: #%d\n", discussion.Number))
	sb.WriteString(fmt.Sprintf("**URL**: %s\n", discussion.URL))
	sb.WriteString(fmt.Sprintf("**Created by**: %s\n\n", discussion.User.Login))
	sb.WriteString(fmt.Sprintf("## Description\n%s\n\n", discussion.Body))

	if len(discussionComments) > 0 {
		sb.WriteString("## Comments\n")
		for i, comment := range discussionComments {
			sb.WriteString(fmt.Sprintf("### Comment %d by %s\n", i+1, comment.User.Login))
			sb.WriteString(fmt.Sprintf("%s\n\n", comment.Body))
		}
	}

	return sb.String()
}
