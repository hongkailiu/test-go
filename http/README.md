# An Web App for SVT Testing

## Build

```bash
$ make build-http

```

## Run

```bash
$ [PORT=8080] ./build/http

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

On Prometheus UI: check on metrics `http_requests_total`, `random_number`, `storage_operation_duration_seconds_count`.