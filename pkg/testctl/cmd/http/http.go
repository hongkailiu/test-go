package http

import (
	"fmt"
	"os"

	"github.com/hongkailiu/test-go/pkg/http"
	"github.com/hongkailiu/test-go/pkg/httpreverse"
	status "github.com/hongkailiu/test-go/pkg/status/server"
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	example = `
# Start http server
testctl http start

# Get http secret
testctl http getSecret`

	log = logrus.New()
)

// NewCmdHTTP ...
func NewCmdHTTP(c *config.Config) *cobra.Command {
	hc := &config.HttpConfig{
		Config:  *c,
		Version: config.VERSION,
	}
	cmd := &cobra.Command{
		Use:     "http",
		Short:   "HTTP server",
		Long:    "HTTP server",
		Example: example,
		Args:    cobra.NoArgs,
	}

	startCmd := &cobra.Command{
		Use:     "start",
		Short:   "Start HTTP server",
		Example: "testctl http start",
		//Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			setup(c, "http")
			http.Run(hc, log)
		},
	}
	startCmd.Flags().BoolVarP(&hc.PProf, "pprof", "p", false, "enable pprof, see https://golang.org/pkg/net/http/pprof/")
	cmd.AddCommand(startCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "getSecret",
		Short: "Get http secret",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(http.GetSecret(32))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Get http status",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			status.Start()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "reverseStart",
		Short: "Start http reverse server",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			setup(c, "reverse")
			httpreverse.Start(log)
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

func setup(c *config.Config, component string) {
	//https://github.com/sirupsen/logrus/pull/653#issuecomment-454467900
	log.SetFormatter(&formatter{
		fields: logrus.Fields{
			"component": component,
		},
		lf: &logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		},
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.WarnLevel)
	if c.Verbose {
		log.SetLevel(logrus.DebugLevel)
	}
}
