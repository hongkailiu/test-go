package ocpsanity

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hongkailiu/test-go/pkg/lib/util"
	"github.com/hongkailiu/test-go/pkg/ocpsanity"
	"github.com/hongkailiu/test-go/pkg/testctl/cmd/config"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	example = `
# Start OCP sanity check
testctl sanity`

	logFilePath = filepath.Join(os.TempDir(), "ocpsanity.log")
	logger      *logrus.Entry
	containLogs bool
)

// NewCmdOCPSanity ...
func NewCmdOCPSanity(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ocpsanity",
		Short:   "OCP sanity check",
		Example: example,
		//Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if containLogs && !c.Verbose {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", "use \"-v\" to show container logs.")
				os.Exit(1)
			}
			logger = newLogger(c)
			logger.WithFields(logrus.Fields{"logFilePath": logFilePath}).Info("logging to file")
			logger.WithFields(logrus.Fields{"startTime": time.Now().Format(time.RFC3339)}).Info("Starting OCP sanity check")
			err := sanityCheck()
			logger.WithFields(logrus.Fields{"endTime": time.Now().Format(time.RFC3339)}).Info("Finishing OCP sanity check")
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
		},
	}
	cmd.Flags().BoolVarP(&containLogs, "containerLogs", "l", false, "show container logs, -v as well")
	return cmd
}

func newLogger(c *config.Config) *logrus.Entry {
	pathMap := lfshook.PathMap{
		logrus.DebugLevel: logFilePath,
		logrus.InfoLevel:  logFilePath,
		logrus.WarnLevel:  logFilePath,
		logrus.ErrorLevel: logFilePath,
		logrus.FatalLevel: logFilePath,
	}

	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logger.SetLevel(logrus.WarnLevel)
	if c.Verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.TextFormatter{},
	))
	return logger.WithFields(logrus.Fields{
		"component": "ocpsanity",
	})
}

func sanityCheck() error {
	return ocpsanity.StartSanityCheck(getConfigPath(), logger, containLogs)
}

func getConfigPath() string {
	configPath := os.Getenv("KUBECONFIG")
	if configPath != "" {
		return configPath
	}

	if home := util.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}

	return ""
}
