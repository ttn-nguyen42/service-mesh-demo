#!/usr/bin/env bash

DEPLOYMENTS=./deployments/consul
cd $DEPLOYMENTS

kubectl create namespace demo

kubectl apply -f locations.yaml \
    -f weather.yaml \
    -f dashboard.yaml

cd -
