package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/hongkailiu/test-go/pkg/version"
)

var Verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "validate-release",
	Short: "A CLI tool to validate OpenShift releases",
	Run: func(cmd *cobra.Command, args []string) {
		if Verbose {
			fmt.Println("Verbose mode is enabled.")
		}
		logrus.WithField("version", version.Version).Info("Start")
		opts, err := GetOptionsFromEnv()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to get options from environment")
		}

		if err := opts.Run(cmd.Context()); err != nil {
			logrus.WithError(err).Fatal("Failed to run command")
		}
		logrus.WithField("version", version.Version).Info("Exit successfully")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Persistent flags are available to this command and all its subcommands.
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "display verbose output")
}
