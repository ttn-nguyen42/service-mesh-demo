#!/usr/bin/env bash

minikube delete -p linkerd
./scripts/minikube/start.sh linkerd

./scripts/helm/linkerd.sh

./scripts/services/linkerd.sh