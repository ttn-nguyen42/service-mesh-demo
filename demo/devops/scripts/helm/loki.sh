#!/usr/bin/env bash

helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
kubectl create namespace prometheus

# Install Loki using Helm
helm install --values ./values/loki.yml \
    loki \
    grafana/loki \
    --namespace prometheus