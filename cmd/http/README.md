# An Web App for SVT Testing

## Build

```bash
$ make build-http

```

## Run

```bash
$ [PORT=8080] ./build/http
$ curl localhost:8080
### this is where the magic goes: prometheus will use this url to obtain metrics data
$ curl localhost:8080/metrics

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
$ buildah bud --format=docker -f test_files/docker/Dockerfile.http.txt -t quay.io/hongkailiu/test-go:http-0.0.1 .
$ buildah push --creds=hongkailiu d58cbf2a06aa docker://quay.io/hongkailiu/test-go:http-0.0.1

```
