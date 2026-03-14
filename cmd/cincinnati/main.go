package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/hongkailiu/test-go/pkg/cincinnati"
)

var (
	port int
)

var rootCmd = &cobra.Command{
	Use:   "cincinnati",
	Short: "A Cincinnati update graph server",
	Long:  "cincinnati is a GoLang implementation of the Red Hat OpenShift Cincinnati update graph protocol",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: get the level from an arg
		logrus.SetLevel(logrus.DebugLevel)
		if err := cincinnati.Start(fmt.Sprintf(":%d", port)); err != nil {
			logrus.WithError(err).Error("Error starting server")
		}
	},
}

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
