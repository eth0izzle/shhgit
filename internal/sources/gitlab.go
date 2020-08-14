package sources

import (
	"context"
	"time"

	"github.com/eth0izzle/shhgit/internal/helpers"
	"github.com/eth0izzle/shhgit/internal/settings"
	"github.com/eth0izzle/shhgit/internal/types"
)

type GitLabProject struct {
	ID          int    `json:"id"`
	CloneURL    string `json:"http_url_to_repo"`
	Stars       int    `json:"star_count"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func StartGitLab(sourceConfig settings.ConfigSource, repositoryChan chan types.RepositoryResource) {
	sleep := time.Duration(sourceConfig.CheckInterval) * time.Second
	ctx := context.Background()
	keys := map[int]bool{}

	for c := time.Tick(sleep); ; {
		var projects []GitLabProject

		if err := helpers.FetchUrlAs(sourceConfig.Endpoint+"&per_page=", "", &projects); err == nil {
			for _, e := range projects {
				if _, value := keys[e.ID]; !value {
					keys[e.ID] = true
					repositoryChan <- types.RepositoryResource{
						Id:          string(e.ID),
						SCM:         types.GIT_SCM,
						Type:        types.GITLAB_SOURCE,
						Location:    e.CloneURL,
						Name:        e.Name,
						Description: e.Description,
						Stars:       e.Stars,
					}
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
