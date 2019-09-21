package quay

type RepositoryEvent struct {
	Name        string   `json:"name"`
	Repository  string   `json:"repository"`
	Namespace   string   `json:"namespace"`
	DockerURL   string   `json:"docker_url"`
	Homepage    string   `json:"homepage"`
	UpdatedTags []string `json:"updated_tags,omitempty"`
}

func (event RepositoryEvent) GetTheMostRecentTag() string {
	return event.UpdatedTags[len(event.UpdatedTags)-1]
}
