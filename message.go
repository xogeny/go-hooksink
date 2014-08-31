package hooksink

type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Repository struct {
	Status    string
	RepoUrl   string `json:"repo_url"`
	Owner     Owner
	IsPrivate bool `json:"is_private"`
	Name      string
	StarCount int    `json:"star_count"`
	RepoName  string `json:"repo_name"`
	GitUrl    string `json:"git_url"`
}

type HeadCommit struct {
	Id string `json:"id"`
}

type PushData struct {
	PushedAt int `json:"pushed_at"`
	Images   []string
	Pusher   string
}

/* A Golang representation of (some of) the GitHub push event payload */
type HubMessage struct {
	Repository Repository `json:"repository"`
	HeadCommit HeadCommit `json:"head_commit"`
	PushData   PushData   `json:"push_data"`

	After string `json:"after"`
}
