package core

import (
	"context"
	"net/url"
	"time"

	"gopkg.in/src-d/go-git.v4"
)

func CloneRepository(session *Session, rawUrl string, dir string) (*git.Repository, error) {
	localCtx, cancel := context.WithTimeout(session.Context, time.Duration(*session.Options.CloneRepositoryTimeout)*time.Second)
	defer cancel()

	if len(session.Config.GitHubEnterpriseUrl) > 0 {
		githubUrl, err := url.Parse(rawUrl)
		if err != nil {
			return nil, err
		}

		userInfo := url.User(session.Config.GitHubAccessTokens[0])
		githubUrl.User = userInfo
		rawUrl = githubUrl.String()
	}

	repository, err := git.PlainCloneContext(localCtx, dir, false, &git.CloneOptions{
		Depth:             1,
		RecurseSubmodules: git.NoRecurseSubmodules,
		URL:               rawUrl,
		SingleBranch:      true,
		Tags:              git.NoTags,
	})

	if err != nil {
		return nil, err
	}

	return repository, nil
}
