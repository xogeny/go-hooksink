package hooksink

/* A Golang representation of (some of) the GitHub push event payload */
type HubMessage struct {
	Repository struct {
		Status    string
		RepoUrl   string `json:"repo_url"`
		Owner     struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
		}
		IsPrivate bool `json:"is_private"`
		Name      string
		StarCount int    `json:"star_count"`
		RepoName  string `json:"repo_name"`
	}

	Push_data struct {
		PushedAt int `json:"pushed_at"`
		Images   []string
		Pusher   string
	}
}
