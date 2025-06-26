#!/usr/bin/env bash
echo "🔧 Initializing Kubernetes cluster..."

kubectl() {
  KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}" command kubectl "$@"
}
function login_to_docker(){
  if kubectl get secret dockerhub-creds -n $1 &>/dev/null; then
    echo "Docker Hub credentials already exist."
    return
  fi
  if [[ -z "$DOCKERHUB_PASSWORD" ]]; then
    echo "Please set the DOCKERBUB_PASSWORD environment variable."
    exit 1
  fi

kubectl create secret -n $1 docker-registry dockerhub-creds \
  --docker-server=https://index.docker.io/v1/ \
  --docker-username=$DOCKERHUB_USER \
  --docker-password=$DOCKERHUB_PASSWORD \
  --docker-email=$EMAIL
}

login_to_docker default
login_to_docker monitoring
echo "Traefik "
bash traefik/init.sh
echo "MetalLB " 
bash metal/init.sh
echo "Cert Manager "
bash certmanager/init.sh
echo "Pi-hole "
bash pihole/init.sh
echo "Graphana "
bash graphana/init.sh
bash graphana/power.sh
