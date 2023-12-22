#!/usr/bin/env bash

# Install Linkerd Helm charts
helm repo add linkerd https://helm.linkerd.io/stable

kubectl create namespace linkerd
kubectl label namespace linkerd pod-security.kubernetes.io/enforce=privileged --overwrite

# Install the CRDs
helm install linkerd-crds linkerd/linkerd-crds \
  -n linkerd --create-namespace --debug

export PATH=$HOME/.linkerd2/bin:$PATH
linkerd install-cni | kubectl apply -f -

# Install the control plane
helm install linkerd-control-plane \
  -n linkerd \
  --set-file identityTrustAnchorsPEM=./certs/ca.crt \
  --set-file identity.issuer.tls.crtPEM=./certs/issuer.crt \
  --set-file identity.issuer.tls.keyPEM=./certs/issuer.key \
  -f ./values/linkerd_cp.yaml \
  linkerd/linkerd-control-plane --debug

# Install the viz extension
helm install linkerd-viz \
  -n linkerd --create-namespace \
   linkerd/linkerd-viz \
  -f ./values/linkerd_viz.yaml --debug

kubectl label namespace linkerd pod-security.kubernetes.io/enforce=privileged --overwrite