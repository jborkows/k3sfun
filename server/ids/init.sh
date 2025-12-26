kubectl() {
  KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}" command kubectl "$@"
}

if [[ -z "$DOMAIN_NAME" ]]; then
  echo "Please set the DOMAIN_NAME environment variable." >&2
  exit 1
fi

if [[ -z "$TZ_TOKEN" ]]; then
  echo "Please set the TX_TOKEN environment variable." >&2
  exit 1
fi
mkdir -p /mnt/external/mainpage

current_script_dir=$(dirname "$(readlink -f "$0")")
(cd "$current_script_dir"

sed "s/DOMAIN/$DOMAIN_NAME/g" ids.yaml |
sed "s/TZ_TOKEN/$TZ_TOKEN/g" > /tmp/ids.yaml
kubectl apply -f /tmp/ids.yaml
)
