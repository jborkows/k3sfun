#!/bin/bash
set -euo pipefail

kubectl() {
  KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}" command kubectl "$@"
}

current_script_dir=$(dirname "$(readlink -f "$0")")

# Check required environment variables
if [[ -z "${DOMAIN_NAME:-}" ]]; then
  echo "Error: Please set the DOMAIN_NAME environment variable."
  exit 1
fi

# Create namespace if needed
if ! kubectl get namespace applications >/dev/null 2>&1; then
  echo "Creating applications namespace..."
  kubectl create namespace applications
else
  echo "Applications namespace already exists."
fi

# Check for OIDC secret
if ! kubectl get secret filebrowser-oidc -n applications >/dev/null 2>&1; then
  echo ""
  echo "=========================================="
  echo "OIDC Secret Not Found"
  echo "=========================================="
  echo ""
  echo "To enable OIDC authentication, create the secret with:"
  echo ""
  echo "  # Generate a random cookie secret (32 bytes, base64 encoded)"
  echo "  COOKIE_SECRET=\$(openssl rand -base64 32 | tr -d '\\n')"
  echo ""
  echo "  kubectl create secret generic filebrowser-oidc \\"
  echo "    --namespace=applications \\"
  echo "    --from-literal=OIDC_CLIENT_ID='your-client-id' \\"
  echo "    --from-literal=OIDC_CLIENT_SECRET='your-client-secret' \\"
  echo "    --from-literal=OIDC_ISSUER_URL='https://your-oidc-provider.com' \\"
  echo "    --from-literal=COOKIE_SECRET=\"\$COOKIE_SECRET\""
  echo ""
  echo "Or use environment variables:"
  echo ""
  echo "  kubectl create secret generic filebrowser-oidc \\"
  echo "    --namespace=applications \\"
  echo "    --from-literal=OIDC_CLIENT_ID=\"\$OIDC_CLIENT_ID\" \\"
  echo "    --from-literal=OIDC_CLIENT_SECRET=\"\$OIDC_CLIENT_SECRET\" \\"
  echo "    --from-literal=OIDC_ISSUER_URL=\"\$OIDC_ISSUER_URL\" \\"
  echo "    --from-literal=COOKIE_SECRET=\"\$(openssl rand -base64 32 | tr -d '\\n')\""
  echo ""
  echo "=========================================="
  echo ""
  
  read -p "Continue deployment without OIDC? (y/N) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Deployment cancelled. Please create the OIDC secret first."
    exit 1
  fi
else
  echo "OIDC secret 'filebrowser-oidc' found."
fi

# Deploy filebrowser
echo "Deploying filebrowser..."
sed "s/DOMAIN/$DOMAIN_NAME/g" "$current_script_dir/filebrowser.yaml" > /tmp/filebrowser.yaml
kubectl apply -f /tmp/filebrowser.yaml

echo ""
echo "Filebrowser deployed successfully!"
echo "URL: https://filebrowser.$DOMAIN_NAME"
