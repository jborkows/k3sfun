#!/usr/bin/env bash

kubectl() {
  KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}" command kubectl "$@"
}
if kubectl get ns monitoring &>/dev/null; then
  echo "Monitoring namespace already exists."
  exit 0
fi
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install kube-prom-stack prometheus-community/kube-prometheus-stack \
  --namespace monitoring --create-namespace
current_script_dir=$(dirname "$(readlink -f "$0")")
if [[ -z "$DOMAIN_NAME" ]]; then
  echo "Please set the DOMAIN_NAME environment variable."
  exit 1
fi
sed "s/DOMAIN/$DOMAIN_NAME/g" $current_script_dir/graphana.yaml > /tmp/graphana.yaml
cat /tmp/graphana.yaml
kubectl apply -f /tmp/graphana.yaml
kubectl apply -f $current_script_dir/power.yaml

