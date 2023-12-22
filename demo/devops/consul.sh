#!/usr/bin/env bash

minikube delete
./scripts/minikube/start.sh

# ./scripts/helm/prom.sh
# ./scripts/helm/grafana.sh
./scripts/helm/consul.sh

./scripts/services/consul.sh