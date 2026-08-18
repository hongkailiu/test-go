package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/hongkailiu/test-go/pkg/cincinnati"
)

var (
	registry string
	repo     string
	output   string
)

var rootCmd = &cobra.Command{
	Use:   "openshift-release-inspect TAG",
	Short: "Inspect an OpenShift release image",
	Long:  "Fetches and displays Cincinnati release metadata from an OpenShift release image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		logrus.SetLevel(logrus.DebugLevel)

		image := args[0]
		result, err := inspect(image)
		if err != nil {
			logrus.WithError(err).WithField("image", image).Fatal("Failed to inspect image")
		}

		var data []byte
		switch output {
		case "json":
			data, err = json.MarshalIndent(result, "", "  ")
		case "yaml":
			data, err = yaml.Marshal(result)
		default:
			return fmt.Errorf("unsupported output format: %s", output)
		}
		if err != nil {
			return fmt.Errorf("failed to marshal output: %w", err)
		}

		if _, err := fmt.Fprintln(os.Stdout, string(data)); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVar(&registry, "registry", "https://quay.io", "Registry URL")
	rootCmd.Flags().StringVar(&repo, "repo", "openshift-release-dev/ocp-release", "Repository in the form org/repo")
	rootCmd.Flags().StringVarP(&output, "output", "o", "yaml", "Output format (yaml or json)")
}

type Result struct {
	cincinnati.ImageInfo
}

func inspect(image string) (Result, error) {
	logrus.WithField("image", image).Debug("Inspecting image")
	var ret Result
	imageInfo, err := cincinnati.GetImageInfo(image)
	if err != nil {
		return ret, err
	}
	ret.ImageInfo = imageInfo
	return ret, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
