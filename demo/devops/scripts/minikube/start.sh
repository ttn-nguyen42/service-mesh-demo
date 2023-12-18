#!/usr/bin/env bash

if minikube status >/dev/null 2>&1; then
    echo "Minikube cluster is running."
else
    echo "Starting Minikube cluster..."
    minikube config set memory 4096
    minikube config set cpus 4
    minikube start
fi

if minikube addons list | grep -q dashboard; then
    echo "Kubernetes dashboard is already installed."
else
    echo "Installing Kubernetes dashboard..."
    minikube addons enable dashboard
fi

if ! kubectl config get-contexts | grep -q minikube; then
    echo "minikube context is not set. Setting minikube context..."
    minikube update-context
fi

echo "Start Kubernetes dashboard"
minikube dashboard
