cincinnati_namespace := env('CINCINNATI_NAMESPACE', "hongkliu-test")
deploy_cincinnati:
    oc get ns "{{cincinnati_namespace}}" || oc create ns "{{cincinnati_namespace}}"
    oc apply -n "{{cincinnati_namespace}}" -f './manifest/*.yaml'


test_cincinnati: deploy_cincinnati
    #!/usr/bin/env bash
    oc -n "{{cincinnati_namespace}}" delete pod -l app=graph-builder --wait
    oc -n "{{cincinnati_namespace}}" scale --replicas=2 deployment graph-builder --timeout=30s
    oc -n "{{cincinnati_namespace}}" wait --timeout=180s --for=condition=Available deployment -l app=graph-builder
    oc -n "{{cincinnati_namespace}}" scale --replicas=3 deployment graph-builder --timeout=30s
    oc -n "{{cincinnati_namespace}}" wait --timeout=180s --for=condition=Available deployment -l app=graph-builder
    set -euxo pipefail
    while read -r pod_name
    do
      oc -n "{{cincinnati_namespace}}" exec "${pod_name}" -- curl -f -s -v "localhost:8080/graph"
      oc -n "{{cincinnati_namespace}}" exec "${pod_name}" -- curl -f -s -v "graph-builder.{{cincinnati_namespace}}.svc.cluster.local/graph"
    done < <(oc -n "{{cincinnati_namespace}}" get pod -l app=graph-builder --no-headers -o custom-columns=":metadata.name")
