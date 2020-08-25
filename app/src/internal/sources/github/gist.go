package github

import (
	"context"
	"fmt"
	"time"

	"github.com/eth0izzle/shhgit/internal/helpers"
	"github.com/eth0izzle/shhgit/internal/settings"
	"github.com/eth0izzle/shhgit/internal/types"
)

func FetchGists(sourceConfig settings.ConfigSource, repositories chan<- types.RepositoryResource) {
	sleep := time.Duration(sourceConfig.CheckInterval) * time.Second
	ctx := context.Background()
	keys := make(map[string]bool)

	for c := time.Tick(sleep); ; {
		var gists []Gist
		token := ""

		if len(sourceConfig.Tokens) > 0 {
			token = fmt.Sprintf("token %s", helpers.GetRandomToken(sourceConfig.Tokens))
		}

		if err := helpers.FetchUrlAs(fmt.Sprintf(sourceConfig.Endpoint+"?per_page=%d", sourceConfig.PerPage), token, &gists); err == nil {
			for _, e := range gists {
				if _, value := keys[e.ID]; !value {
					keys[e.ID] = true

					repositories <- types.RepositoryResource{
						Id:          e.ID,
						SCM:         types.GIT_SCM,
						Type:        types.GIST_SOURCE,
						Location:    e.URL,
						Name:        e.Description,
						Owner:       e.Owner.Username,
						Description: e.Description,
						Stars:       -1,
						Size:        -1,
					}
				}
			}
		} else {
			fmt.Println("Gist request failed - check if your token is valid.." + err.Error())
		}

		select {
		case <-c:
			continue
		case <-ctx.Done():
			return
		}
	}
}
