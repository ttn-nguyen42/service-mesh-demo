#!/usr/bin/env bash

if ! command -v minikube &> /dev/null; then
    echo "minikube is not installed. Installing minikube..."
    curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb
    sudo dpkg -i minikube_latest_amd64.deb
    rm -rf minikube_latest_amd64.deb

fi

if ! command -v kubectl &> /dev/null; then
    echo "kubectl is not installed. Installing kubectl..."
    sudo apt-get update && sudo apt-get install -y apt-transport-https gnupg2
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
    echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
    sudo apt-get update
    sudo apt-get install -y kubectl
fi

if ! command -v helm &> /dev/null; then
    curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get-helm-3 > get_helm.sh
    chmod 700 get_helm.sh
    ./get_helm.sh
    rm -rf get_helm.sh
fi