#!/usr/bin/env bash

set -euo pipefail

cd $(dirname $0)/..

kubeconfig="/tmp/image.pusher.config"

# Make sure user env var wont overrule
unset KUBECONFIG


go build  -v -o /tmp/myhive ./hive
/tmp/myhive \
  --kubeconfig=$kubeconfig