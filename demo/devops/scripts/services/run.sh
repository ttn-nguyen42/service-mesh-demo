#!/usr/bin/env bash

DEPLOYMENTS=./deployments
cd $DEPLOYMENTS

kubectl apply -f .
cd -
