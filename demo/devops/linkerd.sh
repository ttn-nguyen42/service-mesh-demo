#!/usr/bin/env bash

minikube delete
./scripts/minikube/start.sh

./scripts/helm/linkerd.sh

./scripts/services/run.sh