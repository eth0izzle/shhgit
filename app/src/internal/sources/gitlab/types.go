package gitlab

type GitLabProject struct {
	ID          int    `json:"id"`
	CloneURL    string `json:"http_url_to_repo"`
	Stars       int    `json:"star_count"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
