#!/usr/bin/env bash
kubectl() {
  KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}" command kubectl "$@"
}
current_script_dir=$(dirname "$(readlink -f "$0")")
kubectl apply -f $current_script_dir/power.yaml
