# testctl script

## Build

```bash
$ make build-testctl

```

## Run

```bash
$ ./build/testctl

```

## Prometheus

```
# cat prometheus.yml 
...
scrape_configs:
...
  - job_name: 'svt-go'
    static_configs:
    - targets: ['localhost:8080']

# ./prometheus --config.file=prometheus.yml

```

On Prometheus UI: check on metrics `http_requests_total`, `random_number`, `storage_operation_duration_seconds`.

## Postgresql

```bash
# podman run --name some-postgres -p 5432:5432 -e POSTGRES_PASSWORD=redhat -e POSTGRES_USER=redhat -e POSTGRES_DB=ttt -t -i postgres:11.0 

```

## Images

```bash
$ buildah bud --format=docker -f test_files/docker/Dockerfile.testctl.txt -t quay.io/hongkailiu/test-go:testctl-0.0.1 .
$ buildah push --creds=hongkailiu d58cbf2a06aa docker://quay.io/hongkailiu/test-go:testctl-0.0.1

```

## Deployment

```
$ curl https://web-hongkliu-run.b542.starter-us-east-2a.openshiftapps.com

```

Debug deployment

```
$ oc create -f https://raw.githubusercontent.com/hongkailiu/svt-case-doc/master/files/dc_centos.yaml

$ oc get pod -o wide
NAME             READY   STATUS    RESTARTS   AGE   IP              NODE                                          NOMINATED NODE
centos-1-g9kvv   1/1     Running   0          40s   10.128.31.74    ip-172-31-65-156.us-east-2.compute.internal   <none>
web-1-r2smq      1/1     Running   0          8m    10.130.39.153   ip-172-31-69-154.us-east-2.compute.internal   <none>
web-1-zc5kh      1/1     Running   0          8m    10.130.34.186   ip-172-31-74-250.us-east-2.compute.internal   <none>

$ oc rsh centos-1-g9kvv 
sh-4.2$ curl 10.130.39.153:8080
{"version":"0.0.19","ips":["127.0.0.1","::1","10.130.39.153","fe80::1c31:92ff:fecb:42b1"],"now":"2019-03-14T21:28:28.512922733Z"} 

```

Clean up:

```
sh-4.2$ exit
$ oc delete -f https://raw.githubusercontent.com/hongkailiu/svt-case-doc/master/files/dc_centos.yaml

```

Use the binary in the docker image:

```
### interactive shell
# podman run --rm -it -v "/root/.kube/config:/root/.kube/config:ro,z" quay.io/hongkailiu/test-go:testctl-0.0.3-71e8cd9c /bin/ash
/ # ./testctl ocpsanity

### run directly
# podman run --rm -it -v "/root/.kube/config:/root/.kube/config:ro,z" quay.io/hongkailiu/test-go:testctl-0.0.3-71e8cd9c /testctl ocpsanity

```
