package main

import (
    "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    buildv1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
    "path/filepath"
    "flag"
    "k8s.io/client-go/tools/clientcmd"
    "os"
    "fmt"
)

func main() {
    var kubeconfig *string
    if home := homeDir(); home != "" {
        kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    } else {
        kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }
    flag.Parse()

    // use the current context in kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        panic(err.Error())
    }

    buildV1Client, err := buildv1.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    builds, err := buildV1Client.Builds("").List(metav1.ListOptions{})

    if err != nil {
        panic(err.Error())
    }
    fmt.Printf("There are %d builds in the cluster\n", len(builds.Items))

    //Change namespace and build accordingly
    namespace := "testkkk"
    build := "cakephp-mysql-example-1"
    myBuild, err := buildV1Client.Builds(namespace).Get(build, metav1.GetOptions{})
    if errors.IsNotFound(err) {
        fmt.Printf("Build %s in namespace %s not found\n", build, namespace)
    } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
        fmt.Printf("Error getting pod %s in namespace %s: %v\n",
            build, namespace, statusError.ErrStatus.Message)
    } else if err != nil {
        panic(err.Error())
    } else {
        fmt.Printf("Found build %s in namespace %s\n", build, namespace)
        fmt.Printf("Found build %s in status %+v\n", build, myBuild.Status)
    }
}

func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return os.Getenv("USERPROFILE") // windows
}