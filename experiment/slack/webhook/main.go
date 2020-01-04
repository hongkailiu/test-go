package main

import (
	"flag"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

type options struct {
	url string
}

func main() {
	var o options
	flag.StringVar(&o.url, "webhook-url", "", "webhook url")
	flag.Parse()
	if err := slack.PostWebhook(o.url, &slack.WebhookMessage{
		Text:        "Hello6 world!",
		Attachments: nil,
		Channel:     "",
		Parse:       "",
	}); err != nil {
		logrus.WithError(err).Fatal("Failed to send msg to webhook")
	}
}
