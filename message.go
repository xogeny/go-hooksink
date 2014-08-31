package hooksink

// Owner provides information about the owner of a repository
type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Repository contains information about the repository where
// an event took place.
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

// HeadCommit contains information about the commit that triggered
// the event
type HeadCommit struct {
	Id string `json:"id"`
}

// PushData contains information about the push itself
type PushData struct {
	PushedAt int `json:"pushed_at"`
	Images   []string
	Pusher   string
}

// PushMessage is a native Go representation of the (essential) information
// provided with a push event.
//
// Quite a bit of information is missing.  Pull requests filling in the rest
// of the data provided by the GitHub API are welcome.  I just implemented
// what I needed for the moment.
type PushMessage struct {
	// Repository contains information about the repository the push
	// was made to
	Repository Repository `json:"repository"`
	// HeadCommit contains information about the commit that triggered
	// the push
	HeadCommit HeadCommit `json:"head_commit"`
	// PushData contains information about the push itself
	PushData PushData `json:"push_data"`

	// After is the hash of the repository after the commit was completed
	After string `json:"after"`
}
