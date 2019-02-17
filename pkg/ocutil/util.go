package ocutil

import (
	"github.com/hongkailiu/test-go/pkg/lib/util"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	projectv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	"k8s.io/client-go/kubernetes"
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

func (c *CLI) ProjectClient() *projectv1.ProjectV1Client {
	client, err := projectv1.NewForConfig(c.config)
	util.PanicIfError(err)
	return client
}

func (c *CLI) DeployConfigClient() *appsv1.AppsV1Client {
	client, err := appsv1.NewForConfig(c.config)
	util.PanicIfError(err)
	return client
}

func (c *CLI) K8SClientSet() *kubernetes.Clientset {
	clientSet, err := kubernetes.NewForConfig(c.config)
	util.PanicIfError(err)
	return clientSet
}
