package bitbucket

import (
	"context"
	"fmt"
	"time"

	"github.com/eth0izzle/shhgit/internal/helpers"
	"github.com/eth0izzle/shhgit/internal/settings"
	"github.com/eth0izzle/shhgit/internal/types"
)

func Fetch(sourceConfig settings.ConfigSource, repositories chan<- types.RepositoryResource) {
	sleep := time.Duration(sourceConfig.CheckInterval) * time.Second
	ctx := context.Background()
	keys := make(map[string]bool)

	for c := time.Tick(sleep); ; {
		var reposPage BitBucketReposPage

		if err := helpers.FetchUrlAs(fmt.Sprintf(sourceConfig.Endpoint+"?pagelen=%d&after=%s", sourceConfig.PerPage, time.Now().UTC().Add(-sleep).Format(time.RFC3339)), "", &reposPage); err == nil {
			for _, e := range reposPage.Values {
				if _, value := keys[e.ID]; !value {
					scm := types.HG_SCM

					if e.SCM == "git" {
						scm = types.GIT_SCM
					}

					if e.Type == "repository" {
						keys[e.ID] = true
						repositories <- types.RepositoryResource{
							Id:          e.ID,
							Type:        types.BITBUCKET_SOURCE,
							SCM:         scm,
							Location:    e.Links.Clone[0].Href,
							Name:        e.Name,
							Owner:       e.Owner.Username,
							Description: e.Description,
							Size:        e.Size / 1024,
							Stars:       -1,
						}
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
