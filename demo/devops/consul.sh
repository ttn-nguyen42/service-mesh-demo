#!/usr/bin/env bash

minikube delete -p consul
./scripts/minikube/start.sh consul

# ./scripts/helm/prom.sh
# ./scripts/helm/grafana.sh
./scripts/helm/consul.sh

./scripts/services/consul.sh