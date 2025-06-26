#!/usr/bin/env bash
echo "🔧 Installing MetalLB..."
if  kubectl get ns metallb-system &>/dev/null; then
  echo "MetalLB namespace already exists."
  return 0
fi
if [[ -z "$IP_RANGE" ]]; then
  echo "Please set the IP_RANGE environment variable."
  exit 1
fi
  
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.10/config/manifests/metallb-native.yaml
echo "Waiting for MetalLB to be ready..."
until kubectl get pods -n metallb-system | grep -q 'Running'; do
  echo "⏳ Waiting for MetalLB pods to be Running..."
  sleep 2
done

current_script_dir=$(dirname "$(readlink -f "$0")")
sed "s/PORT/$IP_RANGE/g" $current_script_dir/metal.yaml.template | echo
sed "s/PORT/$IP_RANGE/g" $current_script_dir/metal.yaml.template | kubectl apply -f -

