# test-go

* graph-build: leader election via [Lease](https://kubernetes.io/docs/concepts/architecture/leases/). Only master scrapes image registry.
* Graph is cached by the master and can be dumped to slaves via Rest API. Masters become ready faster.
