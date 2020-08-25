package github

import (
	"context"
	"fmt"
	"time"

	"github.com/eth0izzle/shhgit/internal/helpers"
	"github.com/eth0izzle/shhgit/internal/settings"
	"github.com/eth0izzle/shhgit/internal/types"
)

const GITHUB_REPO_URL = "https://www.github.com/%s"

func FetchRepositories(sourceConfig settings.ConfigSource, repositories chan<- types.RepositoryResource) {
	sleep := time.Duration(sourceConfig.CheckInterval) * time.Second
	ctx := context.Background()
	keys := make(map[string]bool)

	for c := time.Tick(sleep); ; {
		var events []GitHubEvent
		token := ""

		if len(sourceConfig.Tokens) > 0 {
			token = fmt.Sprintf("token %s", helpers.GetRandomToken(sourceConfig.Tokens))
		}

		if err := helpers.FetchUrlAs(fmt.Sprintf(sourceConfig.Endpoint+"?per_page=%d", sourceConfig.PerPage), token, &events); err == nil {
			for _, e := range events {
				if _, value := keys[e.ID]; !value {
					keys[e.ID] = true

					if e.ShouldProcess() {
						var repository GitHubRepo

						if err := helpers.FetchUrlAs(fmt.Sprintf(e.Repo.URL), token, &repository); err == nil {
							fromType := types.GITLAB_SOURCE
							description := repository.Description

							if e.Type != "PushEvent" {
								fromType = types.GITHUB_COMMENT
								comment, _ := e.ParsePayload()

								if e.Type == "IssueCommentEvent" {
									description = comment.(*GitHubIssueCommentEvent).Comment.Body
								} else if e.Type == "IssuesEvent" {
									description = comment.(*GitHubIssuesEvent).Issue.Body
								}
							}

							repositories <- types.RepositoryResource{
								Id:          e.ID,
								SCM:         types.GIT_SCM,
								Type:        fromType,
								Location:    fmt.Sprintf(GITHUB_REPO_URL, e.Repo.Name),
								Name:        e.Repo.Name,
								Owner:       e.Repo.Owner.Username,
								Description: description,
								Stars:       repository.Stars,
								Size:        repository.Size,
							}
						}
					}
				}
			}
		} else {
			fmt.Println("GITHUB request failed - check if your token is valid.." + err.Error())
		}

		select {
		case <-c:
			continue
		case <-ctx.Done():
			return
		}
	}
}
