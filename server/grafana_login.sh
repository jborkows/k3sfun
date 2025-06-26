#!/usr/bin/env bash
kubectl() {
  KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}" command kubectl "$@"
}
kubectl get secret -n monitoring kube-prom-stack-grafana -o jsonpath="{.data.admin-password}" | base64 -d && echo
kubectl get secret -n monitoring kube-prom-stack-grafana -o jsonpath="{.data.admin-user}" | base64 -d && echo
