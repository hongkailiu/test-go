# test-go

## Run

```console
$ MOCK_DIR=mock go run ./cmd/cincinnati --log-level debug --max-concurrency 1 --metrics-tls-disabled
```

## Build image

```console
$ podman build -t test001 -f images/cincinnati/Containerfile .
```

## Test metrics endpoint with mTLS

Thanks to David. [Reference](https://github.com/openshift/cluster-version-operator/pull/1271#issuecomment-3715495830).
```console
$ oc --kubeconfig /tmp/ota-stage.c -n openshift-monitoring exec -c prometheus pod/prometheus-k8s-0 -- curl -si  --cacert /etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt --cert /etc/prometheus/secrets/metrics-client-certs/tls.crt --key /etc/prometheus/secrets/metrics-client-certs/tls.key https://cincinnati.cincinnati-go.svc.cluster.local:9090/metrics
HTTP/1.1 200 OK
Content-Type: text/plain; version=0.0.4; charset=utf-8; escaping=underscores
Date: Tue, 31 Mar 2026 02:22:41 GMT
Transfer-Encoding: chunked

# HELP cincinnati_go_request_duration_seconds The HTTP request latencies in seconds.
# TYPE cincinnati_go_request_duration_seconds histogram

$ oc --kubeconfig /tmp/ota-stage.c -n openshift-monitoring exec -c prometheus pod/prometheus-k8s-0 -- curl -si  --cacert /etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt  https://cincinnati.cincinnati-go.svc.cluster.local:9090/metrics
command terminated with exit code 56

$ oc --kubeconfig /tmp/ota-stage.c -n openshift-monitoring exec pod/metrics-server-544445b9-5gwxp -- curl --insecure --cert /etc/tls/metrics-server-client-certs/tls.crt --key /etc/tls/metrics-server-client-certs/tls.key https://cincinnati.cincinnati-go.svc.cluster.local:9090/metrics
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    84  100    84    0     0    641      0 --:--:-- --:--:-- --:--:--   641
unauthorized common name: system:serviceaccount:openshift-monitoring:metrics-server
```

Note that the ServiceMonitor in User Workload Monitoring cannot use `caFile: /etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt` or `scrapeClass: tls-client-certificate-auth` because it works only in Cluster Monitoring.

To make mTLS work, we have to generate a cert/key file for a ServiceAccount which are signed by Kubernetes, just like the monitoring stack does it to `system:serviceaccount:openshift-monitoring:prometheus-k8s`.

## Test metrics endpoint with TLS

Each of the above `curl` commands returns the metrics data.
