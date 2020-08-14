package sources

import (
	"context"
	"fmt"
	"time"

	"github.com/eth0izzle/shhgit/internal/helpers"
	"github.com/eth0izzle/shhgit/internal/settings"
	"github.com/eth0izzle/shhgit/internal/types"
)

const GITHUB_REPO_URL = "https://www.github.com/%s"

type GitHubEvent struct {
	ID   string     `json:"id"`
	Type string     `json:"type"`
	Repo GitHubRepo `json:"repo"`
}

type GitHubRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Stars       int    `json:"stargazers_count"`
	Size        int64  `json:"size"`
}

func StartGitHub(sourceConfig settings.ConfigSource, repositoryChan chan types.RepositoryResource) {
	sleep := time.Duration(sourceConfig.CheckInterval) * time.Second
	ctx := context.Background()
	keys := make(map[string]bool)

	for c := time.Tick(sleep); ; {
		var events []GitHubEvent
		token := ""

		if len(sourceConfig.Tokens) > 0 {
			token = fmt.Sprintf("token %s", helpers.GetRandomToken(sourceConfig))
		}

		if err := helpers.FetchUrlAs(fmt.Sprintf(sourceConfig.Endpoint+"?per_page=%d", sourceConfig.PerPage), token, &events); err == nil {
			for _, e := range events {
				if _, value := keys[e.ID]; !value {
					keys[e.ID] = true

					if e.Type == "PushEvent" {
						var repository GitHubRepo
						if err := helpers.FetchUrlAs(fmt.Sprintf(e.Repo.URL), token, &repository); err == nil {
							repositoryChan <- types.RepositoryResource{
								Id:          e.ID,
								SCM:         types.GIT_SCM,
								Type:        types.GITHUB_SOURCE,
								Location:    fmt.Sprintf(GITHUB_REPO_URL, e.Repo.Name),
								Name:        e.Repo.Name,
								Description: repository.Description,
								Stars:       repository.Stars,
								Size:        repository.Size,
							}
						}
					}
					// } else if *e.Type == "IssueCommentEvent" {
					// 	observedKeys[*e.ID] = true

					// 	dst := &github.IssueCommentEvent{}
					// 	json.Unmarshal(e.GetRawPayload(), dst)
					// 	session.Comments <- *dst.Comment.Body
					// } else if *e.Type == "IssuesEvent" {
					// 	observedKeys[*e.ID] = true

					// 	dst := &github.IssuesEvent{}
					// 	json.Unmarshal(e.GetRawPayload(), dst)
					// 	session.Comments <- dst.Issue.GetBody()
					// }
				}
			}
		}

		select {
		case <-c:
			continue
		case <-ctx.Done():
			return
		}
	}
}

// func GetGists(session *Session) {
// 	localCtx, cancel := context.WithCancel(session.Context)
// 	defer cancel()

// 	observedKeys := map[string]bool{}
// 	opt := &github.GistListOptions{}

// 	var client *GitHubClientWrapper
// 	for c := time.Tick(sleep); ; {
// 		if client != nil {
// 			session.FreeClient(client)
// 		}

// 		client = session.GetClient()
// 		gists, resp, err := client.Gists.ListAll(localCtx, opt)

// 		if err != nil {
// 			if _, ok := err.(*github.RateLimitError); ok {
// 				session.Log.Warn("Token %s[..] rate limited. Reset at %s", client.Token[:10], resp.Rate.Reset)
// 				client.RateLimitedUntil = resp.Rate.Reset.Time
// 				session.FreeClient(client)
// 				break
// 			}

// 			if _, ok := err.(*github.AbuseRateLimitError); ok {
// 				session.Log.Fatal("GitHub API abused detected. Quitting...")
// 			}

// 			session.Log.Warn("Error getting GitHub Gists: %s ... trying again", err)
// 		}

// 		newGists := make([]*github.Gist, 0, len(gists))
// 		for _, e := range gists {
// 			if observedKeys[*e.ID] {
// 				continue
// 			}

// 			newGists = append(newGists, e)
// 		}

// 		for _, e := range newGists {
// 			observedKeys[*e.ID] = true
// 			session.Gists <- *e.GitPullURL
// 		}

// 		opt.Since = time.Now()

// 		select {
// 		case <-c:
// 			continue
// 		case <-localCtx.Done():
// 			cancel()
// 			return
// 		}
// 	}
// }
