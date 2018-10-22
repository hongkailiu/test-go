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

```

