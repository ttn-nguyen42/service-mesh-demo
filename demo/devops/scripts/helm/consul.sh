#!/usr/bin/env bash

# Install Consul onto minikube cluster
helm repo add hashicorp https://helm.releases.hashicorp.com
kubectl create namespace consul
helm install  -f ./values/consul.yaml \
    consul \
    hashicorp/consul \
    --namespace consul --debug \

helm ls --namespace consul
helm status consul --namespace consul