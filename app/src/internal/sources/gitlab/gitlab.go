package gitlab

import (
	"context"
	"time"

	"github.com/eth0izzle/shhgit/internal/helpers"
	"github.com/eth0izzle/shhgit/internal/settings"
	"github.com/eth0izzle/shhgit/internal/types"
)

func Fetch(sourceConfig settings.ConfigSource, repositories chan<- types.RepositoryResource) {
	sleep := time.Duration(sourceConfig.CheckInterval) * time.Second
	ctx := context.Background()
	keys := map[int]bool{}

	for c := time.Tick(sleep); ; {
		var projects []GitLabProject

		if err := helpers.FetchUrlAs(sourceConfig.Endpoint+"&per_page=", "", &projects); err == nil {
			for _, e := range projects {
				if _, value := keys[e.ID]; !value {
					keys[e.ID] = true
					repositories <- types.RepositoryResource{
						Id:          string(e.ID),
						SCM:         types.GIT_SCM,
						Type:        types.GITLAB_SOURCE,
						Location:    e.CloneURL,
						Name:        e.Name,
						Description: e.Description,
						Stars:       e.Stars,
						Size:        -1,
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
