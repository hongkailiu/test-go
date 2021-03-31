package main

import (
	"context"
	"fmt"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	fmt.Println("hello world")
	//var schemes = runtime.NewScheme()
	//utilruntime.Must(hivev1.AddToScheme(schemes))

	if err := hivev1.AddToScheme(scheme.Scheme); err != nil {
		logrus.WithError(err).Fatal("Failed to add hivev1 to scheme")
	}

	c, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		logrus.Fatal("failed to create client")
	}
	ctx := context.TODO()
	clusterPools := &hivev1.ClusterPoolList{}
	listOption := client.MatchingLabels{
		"a": "b",
	}
	if err := c.List(ctx, clusterPools, listOption); err != nil {
		logrus.WithError(err).Fatal("failed to list cluster pools")
	}
	logrus.WithField("size", len(clusterPools.Items)).Info("matching pools")
	for _, pool := range clusterPools.Items {
		logrus.WithField("pool.Name", pool.Name).Info("found the pool")
	}
}
