package version

import (
	"fmt"
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/util"
	"github.com/spf13/cobra"
)

var (
	versionExample = `
		# Print the client and server versions for the current context
		kubectl version`
)

type Options struct {
	Short      bool
}

func NewCmdVersion() *cobra.Command {
	o := &Options{}
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the client and server version information",
		Long:    "Print the client and server version information for the current context",
		Example: versionExample,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Complete(cmd))
			util.CheckErr(o.Validate())
			util.CheckErr(o.Run())
		},
	}
	cmd.Flags().BoolVar(&o.Short, "short", o.Short, "Print just the version number.")
	return cmd
}

func (o *Options) Complete(cmd *cobra.Command) error {
	return nil
}

func (o *Options) Validate() error {

	return nil
}

func (o *Options) Run() error {
	fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	return nil
}