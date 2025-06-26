#!/usr/bin/env bash

if kubectl get ns cert-manager &>/dev/null; then
  echo "cert-manager namespace already exists."
  exit 0
fi

if [[ -z "$CLOUDFLARE_API_TOKEN" ]]; then
  echo "Please set the CLOUDFLARE_API_TOKEN environment variable."
  exit 1
fi
if [[ -z "$EMAIL" ]]; then
  echo "Please set the EMAIL environment variable."
  exit 1
fi
if [[ -z "$DOMAIN_NAME" ]]; then
  echo "Please set the DOMAIN_NAME environment variable."
  exit 1
fi
if [[ -z "$CERTIFICATE_NAME" ]]; then
  echo "Please set the CERTIFICATE_NAME environment variable."
  exit 1
fi
set -euo pipefail

echo "🔧 Installing cert-manager..."
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/latest/download/cert-manager.yaml
  until kubectl get pods -n cert-manager | grep -q 'Running'; do
    echo "⏳ Waiting for cert-manager pods to be Running..."
    sleep 2
done
kubectl create secret generic cloudflare-api-token-secret \
  --from-literal=api-token="$CLOUDFLARE_API_TOKEN" \
  -n cert-manager

current_script_dir=$(dirname "$(readlink -f "$0")")
sed "s/DOMAIN_NAME/$DOMAIN_NAME/g" $current_script_dir/certmanager.yaml | \
	sed "s/EMAIL/$EMAIL/g" | \
	sed "s/CERTNAME/$CERTIFICATE_NAME/g" | \
	sed "s/DOMAIN/$DOMAIN_NAME/g" > /tmp/certmanager.yaml
cat /tmp/certmanager.yaml
kubectl apply -f /tmp/certmanager.yaml
