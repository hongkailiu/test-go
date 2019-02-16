package ocpsanity

import (
	"github.com/hongkailiu/test-go/pkg/ocutil"
	log "github.com/sirupsen/logrus"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	oc *ocutil.CLI
)

func StartSanityCheck(configPath string) error {
	oc = ocutil.NewCLI(configPath)
	projectList, err := oc.ProjectClient().Projects().List(meta.ListOptions{})
	if err != nil {
		return err
	}
	for _, project := range projectList.Items {
		log.WithFields(log.Fields{"name": project.Name}).Info("Handle project")
	}
	return nil
}
