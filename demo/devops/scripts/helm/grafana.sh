#/usr/bin/env bash

# Install Grafana for Consul
helm repo add grafana https://grafana.github.io/helm-charts
kubectl create namespace prometheus
helm install --values ./values/grafana.yaml \
    grafana \
    grafana/grafana \
    --namespace prometheus