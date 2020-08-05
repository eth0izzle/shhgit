package core

import (
	"context"
	"net/url"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
)

type GitResourceType int

const (
	LOCAL_SOURCE GitResourceType = iota
	GITHUB_SOURCE
	GIST_SOURCE
	BITBUCKET_SOURCE
	GITLAB_SOURCE
)

type GitResource struct {
	Id   int64
	Type GitResourceType
	Url  string
}

func CloneRepository(session *Session, rawUrl string, dir string) (*git.Repository, error) {
	timeout := time.Duration(*session.Options.CloneRepositoryTimeout) * time.Second
	localCtx, cancel := context.WithTimeout(context.Background(), timeout)
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

	session.Log.Debug("[%s] Cloning in to %s", rawUrl, strings.Replace(dir, *session.Options.TempDirectory, "", -1))
	repository, err := git.PlainCloneContext(localCtx, dir, false, &git.CloneOptions{
		Depth:             1,
		RecurseSubmodules: git.NoRecurseSubmodules,
		URL:               rawUrl,
		SingleBranch:      true,
		Tags:              git.NoTags,
	})

	if err != nil {
		session.Log.Debug("[%s] Cloning failed: %s", rawUrl, err.Error())
		return nil, err
	}

	return repository, nil
}
