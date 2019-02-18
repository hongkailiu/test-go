# ocpsanity

StartSanityCheck is a script to check sanity of OCP installation. It requires a good _k8s config file_.
By default, it uses `${HOME}/.kube/config`.
It can be overwritten by the environment variable `KUBECONFIG`, eg, `KUBECONFIG=/root/20190215_124136/auth/kubeconfig`.


## build

```bash
$ make build-ocpsanity

```

## run

```bash
$ ./build/ocpsanity

```

use release version:

```
$ curl -LO https://github.com/hongkailiu/test-go/releases/download/0.0.8/ocpsanity-0.0.8-Linux-x86_64.tar.gz
$ tar xzvf ocpsanity-0.0.8-Linux-x86_64.tar.gz
$ ./ocpsanity/build/ocpsanity

```