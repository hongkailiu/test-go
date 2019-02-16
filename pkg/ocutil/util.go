package ocutil

import (
	"github.com/hongkailiu/test-go/pkg/lib/util"
	project "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type CLI struct {
	configPath string
	config     *restclient.Config
}

func NewCLI(configPath string) *CLI {
	client := &CLI{}
	client.configPath = configPath
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	util.PanicIfError(err)
	client.config = config
	return client
}

func (c *CLI) ProjectClient() *project.ProjectV1Client {
	client, err := project.NewForConfig(c.config)
	util.PanicIfError(err)
	return client
}
