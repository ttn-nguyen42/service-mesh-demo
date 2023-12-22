#!/usr/bin/env bash

DEPLOYMENTS=./deployments/linkerd
cd $DEPLOYMENTS

kubectl create namespace demo

kubectl apply -f locations.yaml \
    -f weather.yaml \
    -f dashboard.yaml

cd -
