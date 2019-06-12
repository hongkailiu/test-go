package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionExample = `
		# Print the client and server versions for the current context
		kubectl version`
)

type Options struct {
	Short bool
}

func NewCmdVersion() *cobra.Command {
	o := &Options{}
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the client and server version information",
		Long:    "Print the client and server version information for the current context",
		Example: versionExample,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		},
	}
	cmd.Flags().BoolVar(&o.Short, "short", o.Short, "Print just the version number.")
	return cmd
}
