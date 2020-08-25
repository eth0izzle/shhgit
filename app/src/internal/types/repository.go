package types

import (
	"fmt"
)

type RepositoryResourceType int
type RepositorySCMType int

const (
	LOCAL_SOURCE RepositoryResourceType = iota
	GITHUB_SOURCE
	GITHUB_COMMENT
	GIST_SOURCE
	BITBUCKET_SOURCE
	GITLAB_SOURCE
)

const (
	GIT_SCM RepositorySCMType = iota
	HG_SCM
)

type RepositoryResource struct {
	Id          string
	Name        string
	Description string
	Owner       string
	SCM         RepositorySCMType
	Type        RepositoryResourceType
	Size        int64
	Location    string
	Stars       int
}

func (e RepositoryResourceType) String() string {
	switch e {
	case LOCAL_SOURCE:
		return "Local"
	case GITHUB_SOURCE:
		return "GitHub"
	case GITHUB_COMMENT:
		return "GitHub Comment"
	case GIST_SOURCE:
		return "Gist"
	case BITBUCKET_SOURCE:
		return "BitBucket"
	case GITLAB_SOURCE:
		return "GitLab"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

func (e RepositoryResource) Handle() error {
	fmt.Println("handling " + e.Location)
	return nil
}
