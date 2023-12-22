#!/usr/bin/env bash

CERT_PATH="./certs"

# Generate mTLS certificates for Linkerd

# Install step CLI
if step >/dev/null 2>&1; then
else
    wget https://dl.smallstep.com/cli/docs-cli-install/latest/step-cli_amd64.deb
    sudo dpkg -i step-cli_amd64.deb
    rm -rf step-cli_amd64.deb
fi

# Trust anchor certificates
step certificate create root.linkerd.cluster.local \
    ca.crt ca.key \
    --profile root-ca \
    --no-password \
    --insecure

mv ca.crt ca.key $CERT_PATH

# Issuer certificate and key
step certificate create identity.linkerd.cluster.local \
    issuer.crt issuer.key \
    --profile intermediate-ca \
    --not-after 8760h \
    --no-password \
    --insecure \
    --ca ca.crt --ca-key ca.key

mv issuer.crt issuer.key $CERT_PATH