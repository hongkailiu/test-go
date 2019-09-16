package github

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type PingEvent struct {
	Zen    string `json:"zen"`
	HookID int    `json:"hook_id"`
	Hook   Hook   `json:"hook"`
}

type Hook struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// TODO ... other fields
}

type Issue struct {
	Number int `json:"number"`
}

type IssueComment struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type IssueCommentEvent struct {
	Action  string       `json:"action"`
	Comment IssueComment `json:"comment"`
}

type WebhookHeaders struct {
	EventType     string `header:"X-GitHub-Event"`
	GUID          string `header:"X-GitHub-Delivery"`
	HMACHexDigest string `header:"X-Hub-Signature"`
	ContentType   string `header:"content-type"`
}

func Handle(payload []byte, wh WebhookHeaders, log *logrus.Logger) (int, string) {
	if wh.EventType == "" {
		return http.StatusBadRequest, "Missing X-GitHub-Event Header"
	}
	if wh.GUID == "" {
		return http.StatusBadRequest, "Missing X-GitHub-Delivery Header"
	}
	if wh.HMACHexDigest == "" {
		return http.StatusForbidden, "Missing X-Hub-Signature"
	}
	if wh.ContentType != "application/json" {
		return http.StatusBadRequest, "Hook only accepts content-type: application/json - please reconfigure this hook on GitHub"
	}

	if !validatePayloadWithSecret(payload, wh.HMACHexDigest) {
		return http.StatusForbidden, "403 Forbidden: Invalid X-Hub-Signature"
	}

	switch wh.EventType {
	case "ping":
		go handlePing(payload, log)
	case "issue_comment":
		go handleIssueComments(payload, log)
	default:
		log.WithField("wh.EventType", wh.EventType).Info("Ignored payload with an unimplemented event type.")
	}
	return http.StatusOK, "OK"
}

func validatePayloadWithSecret(payload []byte, secret string) bool {
	return true
}

func handlePing(payload []byte, log *logrus.Logger) {
	var event PingEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		log.WithError(err).WithField("payload", string(payload)).Error("Cannot parse payload")
		return
	}
	log.WithField("event.HookID", event.HookID).Debug("received PingEvent")
}

func handleIssueComments(payload []byte, log *logrus.Logger) {
	var event IssueCommentEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		log.WithError(err).WithField("payload", string(payload)).Error("Cannot parse payload")
		return
	}
	log.WithField("event.Comment.Body", event.Comment.Body).Debug("received IssueCommentEvent")
}
