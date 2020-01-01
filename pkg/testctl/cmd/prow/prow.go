package prow

import (
	"github.com/hongkailiu/test-go/pkg/prow"
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	example = `
# Monitor Prow deployment
testctl prow monitor`
)

func NewCmdProw(c *config.Config) *cobra.Command {
	pc := &config.ProwConfig{
		Config: *c,
	}

	cmd := &cobra.Command{
		Use:     "prow",
		Short:   "Prow",
		Long:    "Prow",
		Example: example,
		Args:    cobra.NoArgs,
	}

	monitorCmd := &cobra.Command{
		Use:     "monitor",
		Short:   "Monitor Prow deployment",
		Example: "testctl prow monitor",
		Run: func(cmd *cobra.Command, args []string) {
			if err := prow.Monitor(pc); err != nil {
				logrus.WithError(err).Fatalf("Failed to monitor prow!")
			}
		},
	}
	monitorCmd.Flags().StringVar(&pc.KubeConfigPath, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests.")
	cmd.AddCommand(monitorCmd)
	return cmd
}
