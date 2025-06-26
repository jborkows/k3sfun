#!/usr/bin/env bash
set -euo pipefail
EMAIL=$1 
CF_TOKEN=$2
IP_RANGE=$3
CLOUDFLARE_API_TOKEN=$4
if [[ -z "$PASS" ]]; then
  echo "Usage: $0 <ssh_password>"
  exit 1
fi

function install_utils(){
  echo "🔧 Installing utilities..."
  apt-get update
  apt-get install -y curl wget jq git
  mkdir -p /data || true
  mkdir -p /root/yaml || true

  if ! command -v yq &>/dev/null; then
    echo "yq not found, installing yq..."
    wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64
    chmod +x /usr/local/bin/yq
  fi

  sleep_config=/etc/systemd/sleep.conf.d/nosuspend.conf
  if [ ! -d "$(dirname "$sleep_config")" ]; then
    echo "Creating directory for sleep config..."
    mkdir -p "$(dirname "$sleep_config")"
  fi
  if [ ! -f "$sleep_config" ]; then
    echo "Creating sleep config to disable suspend..."
    echo "[Sleep]" >> "$sleep_config" 
    echo "AllowSuspend=no" >> "$sleep_config" 
    echo "AllowHibernation=no" >>"$sleep_config" 
    echo "AllowHybridSleep=no" >>"$sleep_config" 
  else
    echo "Sleep config already exists, skipping creation."
  fi
}

function check_k3s(){
  if command -v k3s >/dev/null 2>&1; then
    echo "k3s binary found at: $(which k3s)"
    echo "k3s version: $(k3s --version)"
    return 0
  fi
    echo "k3s binary not found."
    echo "$PASS" | sudo -S bash -c 'curl -sfL https://get.k3s.io | sh -'
  echo "✅ K3s cluster installed!"
}

function wait_on_k3s(){
  # Wait for node to be Ready
  until echo "$PASS" | sudo -S kubectl get node 2>/dev/null | grep -q ' Ready '; do
    echo "⏳ Waiting for K3s node to be ready..."
    sleep 2
  done

  # Wait for core pods to be Running
  until echo "$PASS" | sudo -S kubectl get pods -n kube-system 2>/dev/null | grep -v Completed | grep -v Evicted | grep -v Running; do
    echo "⏳ Waiting for core system pods..."
    sleep 2
  done

  echo "✅ K3s cluster is fully ready!"
}



echo "$PASS" | sudo -S -i
echo "$PASS" | sudo -S bash -c "$(declare -f install_utils);   install_utils"
check_k3s
wait_on_k3s
git clone https://github.com/jborkows/k3sfun.git
echo "PLEASE LOG INTO SERVER AND RUN ./k3sfun/init.sh"
