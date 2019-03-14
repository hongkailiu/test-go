package cmd

import (
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/flags"
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/http"
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/version"
	"github.com/spf13/cobra"
)

const (
	VERSION = "0.0.1"
)

// NewDefaultTestctlCommand creates the `testctl` command with default arguments
func NewDefaultTestctlCommand() *cobra.Command {
	f := &flags.Flags{}
	cmd := &cobra.Command{
		Use:   "testctl",
		Short: "testctl controlls test program",
		Long: `testctl controlls test program.
	Find more information at https://github.com/hongkailiu/test-go`,
		Run:     runHelp,
		Version: VERSION,
		//BashCompletionFunction: bashCompletionFunc,
	}
	cmd.PersistentFlags().BoolVarP(&f.Verbose, "verbose", "v", false, "verbose output")
	cmd.AddCommand(version.NewCmdVersion())

	cmd.AddCommand(http.NewCmdHTTP(f))

	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
