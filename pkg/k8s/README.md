# Use [k8s client-go](https://github.com/kubernetes/client-go)


## [Install kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/): Optional

```sh
$ curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.9.0/bin/linux/amd64/kubectl
$ chmod +x ./kubectl
```

Put the command line binary <code>kubectl</code> into <code>${PATH}</code>.

```sh
### DID oc-login already
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"9", GitVersion:"v1.9.0", GitCommit:"925c127ec6b946659ad0fd596fa959be43f0cc05", GitTreeState:"clean", BuildDate:"2017-12-15T21:07:38Z", GoVersion:"go1.9.2", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"6", GitVersion:"v1.6.1+5115d708d7", GitCommit:"fff65cf", GitTreeState:"clean", BuildDate:"2017-12-19T23:46:54Z", GoVersion:"go1.7.6", Compiler:"gc", Platform:"linux/amd64"}

```

## Build & run

```sh
$ make build-k8s
$ ./build/k8s
```