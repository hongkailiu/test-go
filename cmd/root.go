package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/hongkailiu/test-go/pkg/version"
)

var rootCmd = &cobra.Command{
	Use:   "validate-release",
	Short: "A CLI tool to validate OpenShift releases",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.WithField("version", version.Version).Info("Start")
		opts, err := GetOptionsFromEnv()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to get options from environment")
		}

		if err := opts.Run(cmd.Context()); err != nil {
			logrus.WithError(err).Fatal("Failed to run command")
		}
		logrus.Info("Exit successfully")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
