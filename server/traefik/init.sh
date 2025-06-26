MANIFEST="/var/lib/rancher/k3s/server/manifests/traefik.yaml"
BACKUP="${MANIFEST}.bak.$(date +%F_%T)"
NAMESPACE="kube-system"
SECRET_NAME="cloudflare-secret"


if  kubectl  get secret "$SECRET_NAME" &>/dev/null; then
	echo "Kubernetes secret $SECRET_NAME already exists"
	exit 0
fi


if [[ -z "${EMAIL:-}" ]]; then
  echo "Please set the EMAIL environment variable."
  exit 1
fi

if [[ -z "${CF_TOKEN:-}" ]]; then
  echo "Please set the CF_TOKEN environment variable with your Cloudflare API token."
  exit 1
fi

if ! kubectl -n "$NAMESPACE" get secret "$SECRET_NAME" &>/dev/null; then
  kubectl -n "$NAMESPACE" create secret generic "$SECRET_NAME" --from-literal=api-token="$CF_TOKEN"
  echo "Created Kubernetes secret $SECRET_NAME"
else
  echo "Kubernetes secret $SECRET_NAME already exists"
fi
if ! kubectl  get secret "$SECRET_NAME" &>/dev/null; then
  kubectl  create secret generic "$SECRET_NAME" --from-literal=api-token="$CF_TOKEN"
  echo "Created Kubernetes secret $SECRET_NAME"
else
  echo "Kubernetes secret $SECRET_NAME already exists"
fi

curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
current_script_dir=$(dirname "$(readlink -f "$0")")
echo "Current script directory: $current_script_dir"
cp ${current_script_dir}/traefik-config.yaml /var/lib/rancher/k3s/server/manifests/traefik-config.yaml
