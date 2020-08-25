package bitbucket

type BitBucketReposPage struct {
	Values []BitBucketRepo `json:"values"`
}

type BitBucketRepo struct {
	Links       BitBucketRepoLinks `json:"links"`
	SCM         string             `json:"scm"`
	Size        int64              `json:"size"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Type        string             `json:"type"`
	ID          string             `json:"uuid"`
	Owner       BitBucketOwner     `json:"owner"`
}

type BitBucketRepoLinks struct {
	Clone []struct {
		Name string `json:"name"`
		Href string `json:"href"`
	} `json:"clone"`
}

type BitBucketOwner struct {
	Username string `json:"display_name"`
}
