package assigner

import "time"

type Status struct {
	Current []string `json:"current"`
	Config  `json:"config,inline"`
}

type Action struct {
	At      time.Time `json:"at"`
	Members []string `json:"members"`
}

type Config struct {
	GroupName        string  `json:"groupName"`
	ScheduledActions []Action `json:"scheduledActions"`
}
