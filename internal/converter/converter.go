package converter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bigwhite/issue2md/internal/github"
)

// Helper to write reactions for a comment or discussion comment.
func writeReactions(sb *strings.Builder, reactions []github.Reaction) {
	if len(reactions) == 0 {
		return
	}
	sb.WriteString("**Reactions:**\n")
	reactionMap := make(map[string][]github.User)
	for _, reaction := range reactions {
		reactionMap[reaction.Content] = append(reactionMap[reaction.Content], reaction.User)
	}
	// Sort reaction types for deterministic output
	reactionTypes := make([]string, 0, len(reactionMap))
	for reaction := range reactionMap {
		reactionTypes = append(reactionTypes, reaction)
	}
	sort.Strings(reactionTypes)
	for _, reaction := range reactionTypes {
		users := reactionMap[reaction]
		sb.WriteString(fmt.Sprintf("- :%s: by %d user(s):\n", reaction, len(users)))
		for _, user := range users {
			sb.WriteString(fmt.Sprintf("  - [%s](https://github.com/%s)\n", user.Login, user.Login))
		}
	}
	sb.WriteString("\n")
}

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
			writeReactions(&sb, comment.Reactions)
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
			writeReactions(&sb, comment.Reactions)
		}
	}

	return sb.String()
}
