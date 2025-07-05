kubectl() {
  KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}" command kubectl "$@"
}
if kubectl get deployment "filebrowser" -n "applications" -o name >/dev/null 2>&1; then
  echo "File browser deployed."
  exit 0
fi
if ! kubectl get namespace applications >/dev/null 2>&1; then
  echo "Creating applications namespace..."
  kubectl create namespace applications
else
  echo "Filebrowser namespace already exists."
fi
current_script_dir=$(dirname "$(readlink -f "$0")")
if [[ -z "$DOMAIN_NAME" ]]; then
  echo "Please set the DOMAIN_NAME environment variable."
  exit 1
fi
cp $current_script_dir/externaldisk.yaml /tmp/externaldisk.yaml 
sed "s/DOMAIN/$DOMAIN_NAME/g" $current_script_dir/filebrowser.yaml > /tmp/filebrowser.yaml
kubectl apply -f /tmp/filebrowser.yaml
