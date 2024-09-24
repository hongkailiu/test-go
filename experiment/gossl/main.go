package main

import (
	"github.com/hongkailiu/test-go/pkg/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := cmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Failed to execute the cmd")
	}
}
