package cmd

import (
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/config"
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/http"
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/ocpsanity"
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/version"
	"github.com/spf13/cobra"
)

// NewDefaultTestctlCommand creates the `testctl` command with default arguments
func NewDefaultTestctlCommand() *cobra.Command {
	c := &config.Config{}
	cmd := &cobra.Command{
		Use:   "testctl",
		Short: "testctl controlls test program",
		Long: `testctl controlls test program.
	Find more information at https://github.com/hongkailiu/test-go`,
		Run:     runHelp,
		Version: config.VERSION,
		//BashCompletionFunction: bashCompletionFunc,
	}
	cmd.PersistentFlags().BoolVarP(&c.Verbose, "verbose", "v", false, "verbose output")
	cmd.AddCommand(version.NewCmdVersion())

	cmd.AddCommand(http.NewCmdHTTP(c))
	cmd.AddCommand(ocpsanity.NewCmdOCPSanity(c))

	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
