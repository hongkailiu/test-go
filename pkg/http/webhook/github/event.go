package github

type Event struct {
	PingEvent
}

type PingEvent struct {
	Zen    string `json:"zen"`
	HookID int    `json:"hook_id"`
	Hook   Hook   `json:"hook"`
}

type Hook struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WebhookHeader struct {
	EventName     string `header:"X-GitHub-Event"`
	GUID          string `header:"X-GitHub-Delivery"`
	HMACHexDigest string `header:"X-Hub-Signature"`
}
