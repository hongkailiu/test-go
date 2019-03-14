package http

import (
	"fmt"
	"os"

	"github.com/hongkailiu/test-go/pkg/http"
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/flags"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	example = `
# Start http server
testctl http start

# Get http secret
testctl http getSecret`
)

func NewCmdHTTP(f *flags.Flags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "http",
		Short:   "HTTP server",
		Long:    "HTTP server",
		Example: example,
		Args:    cobra.NoArgs,
	}

	cmd.AddCommand(&cobra.Command{
		Use:     "start",
		Short:   "Start HTTP server",
		Example: "testctl http start",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			setup(f)
			http.Run()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "getSecret",
		Short: "Get http secret",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(http.GetSecret(32))
		},
	})
	return cmd
}

// formatter adds default fields to each log entry.
type formatter struct {
	fields logrus.Fields
	lf     logrus.Formatter
}

// Format satisfies the logrus.Formatter interface.
func (f *formatter) Format(e *logrus.Entry) ([]byte, error) {
	for k, v := range f.fields {
		e.Data[k] = v
	}
	return f.lf.Format(e)
}

func setup(f *flags.Flags) {
	//https://github.com/sirupsen/logrus/pull/653#issuecomment-454467900
	logrus.SetFormatter(&formatter{
		fields: logrus.Fields{
			"component": "http",
		},
		lf: &logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		},
	})

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)
	if f.Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
