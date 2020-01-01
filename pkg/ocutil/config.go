package ocutil

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func LoadKubeConfigs(kubeconfig string) (map[string]rest.Config, string, error) {
	loader := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	cfg, err := loader.Load()
	if err != nil {
		return nil, "", err
	}
	configs := map[string]rest.Config{}
	for context := range cfg.Contexts {
		contextCfg, err := clientcmd.NewNonInteractiveClientConfig(*cfg, context, &clientcmd.ConfigOverrides{}, loader).ClientConfig()
		if err != nil {
			return nil, "", fmt.Errorf("create %s client: %v", context, err)
		}
		configs[context] = *contextCfg
		logrus.Debugf("Parsed kubeconfig context: %s", context)
	}
	return configs, cfg.CurrentContext, nil
}
