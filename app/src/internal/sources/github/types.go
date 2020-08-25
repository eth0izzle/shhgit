package github

import (
	"encoding/json"

	"github.com/eth0izzle/shhgit/internal/helpers"
)

type GitHubEvent struct {
	ID         string          `json:"id,omitempty"`
	Type       string          `json:"type,omitempty"`
	Repo       GitHubRepo      `json:"repo,omitempty"`
	RawPayload json.RawMessage `json:"payload,omitempty"`
}

type GitHubRepo struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	URL         string      `json:"url,omitempty"`
	Stars       int         `json:"stargazers_count,omitempty"`
	Size        int64       `json:"size,omitempty"`
	Owner       GitHubOwner `json:"owner,omitempty"`
}

type GitHubOwner struct {
	Username string `json:"login,omitempty"`
}

type GitHubIssueCommentEvent struct {
	Comment GitHubComment `json:"comment,omitempty"`
}

type GitHubIssuesEvent struct {
	Issue GitHubIssue `json:"issue,omitempty"`
}

type GitHubComment struct {
	Body string `json:"body,omitempty"`
}

type GitHubIssue struct {
	Body string `json:"body,omitempty"`
}

type Gist struct {
	ID          string      `json:"id,omitempty"`
	URL         string      `json:"git_pull_url,omitempty"`
	Description string      `json:"description,omitempty"`
	Owner       GitHubOwner `json:"owner,omitempty"`
}

func (e *GitHubEvent) ShouldProcess() bool {
	return helpers.Contains([]string{"`PushEvent`", "IssueCommentEvent", "IssuesEvent"}, e.Type)
}

func (e *GitHubEvent) ParsePayload() (payload interface{}, err error) {
	switch e.Type {
	case "IssueCommentEvent":
		payload = &GitHubIssueCommentEvent{}
	case "IssuesEvent":
		payload = &GitHubIssuesEvent{}
	}

	err = json.Unmarshal(e.RawPayload, &payload)
	return payload, err
}
