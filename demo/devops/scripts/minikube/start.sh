#!/usr/bin/env bash

PROFILE=$1

if minikube status -p $PROFILE >/dev/null 2>&1; then
    echo "Minikube cluster is running."
else
    echo "Starting Minikube cluster..."
    minikube config set memory 3072
    minikube config set cpus 2
    minikube start -p $PROFILE
fi

if minikube addons list -p $PROFILE | grep -q dashboard; then
    echo "Kubernetes dashboard is already installed."
else
    echo "Installing Kubernetes dashboard..."
    minikube addons enable metrics-server -p $PROFILE
    minikube addons enable dashboard -p $PROFILE
fi

if ! kubectl config get-contexts | grep -q minikube; then
    echo "minikube context is not set. Setting minikube context..."
    minikube update-context -p $PROFILE
fi

echo "Start Kubernetes dashboard"
minikube dashboard --url -p $PROFILE &
