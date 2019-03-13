package cmd

import (
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/version"
	"github.com/spf13/cobra"
)

// NewDefaultTestctlCommand creates the `testctl` command with default arguments
func NewDefaultTestctlCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "hugo",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	cmds.AddCommand(version.NewCmdVersion())

	return cmds
}
