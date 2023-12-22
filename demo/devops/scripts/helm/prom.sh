#!/usr/bin/env bash

# Install Prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
kubectl create namespace prometheus
helm install --values ./values/prom.yml \
    prometheus \
    prometheus-community/prometheus \
    --namespace prometheus