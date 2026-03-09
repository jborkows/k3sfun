#!/bin/bash
set -euo pipefail

# Vikunja Deployment Script
# Usage: export DOMAIN_NAME=your-domain.com && ./init.sh

kubectl() {
  KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}" command kubectl "$@"
}

current_script_dir=$(dirname "$(readlink -f "$0")")

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo_info() {
  echo -e "${GREEN}[INFO]${NC} $1"
}

echo_warn() {
  echo -e "${YELLOW}[WARN]${NC} $1"
}

echo_error() {
  echo -e "${RED}[ERROR]${NC} $1"
}

# Check required environment variables
if [[ -z "${DOMAIN_NAME:-}" ]]; then
  echo_error "Please set the DOMAIN_NAME environment variable."
  echo "Example: export DOMAIN_NAME=example.com"
  exit 1
fi

echo_info "Deploying Vikunja to k3s cluster..."
echo_info "Domain: $DOMAIN_NAME"

# Create namespace
echo_info "Creating namespace..."
if ! kubectl get namespace vikunja >/dev/null 2>&1; then
  kubectl create namespace vikunja
  echo_info "Namespace 'vikunja' created"
else
  echo_info "Namespace 'vikunja' already exists"
fi

# Ensure external directories exist
echo_info "Checking external storage..."
if [[ ! -d "/mnt/external" ]]; then
  echo_warn "/mnt/external not found. Please ensure external disk is mounted."
  echo "Vikunja will create directories automatically, but data will be lost if not persisted."
fi

# Check for OIDC secret
echo_info "Checking OIDC configuration..."
if ! kubectl get secret vikunja-oidc -n vikunja >/dev/null 2>&1; then
  echo_warn "OIDC secret 'vikunja-oidc' not found in namespace 'vikunja'"
  echo ""
  echo "To enable OIDC authentication, create the secret with:"
  echo ""
  echo "  # Generate a random cookie secret"
  echo "  COOKIE_SECRET=\$(openssl rand -base64 32 | tr -d '\\n')"
  echo ""
  echo "  kubectl create secret generic vikunja-oidc \\"
  echo "    --namespace=vikunja \\"
  echo "    --from-literal=OIDC_CLIENT_ID='<your-tsidp-client-id>' \\"
  echo "    --from-literal=OIDC_CLIENT_SECRET='<your-tsidp-secret>' \\"
  echo "    --from-literal=OIDC_ISSUER_URL='https://idp-1.tailf15a72.ts.net/' \\"
  echo "    --from-literal=CO\"\$COOKIE_SECRET\""
  echo ""
  echo "Or use your existing values:"
  echo ""
  echo "  kubectl create secret generic vikunja-oidc \\"
  echo "    --namespace=vikunja \\"
  echo "    --from-literal=OIDC_CLIENT_ID=\"\$VIKUNJA_OIDC_CLIENT_ID\" \\"
  echo "    --from-literal=OIDC_CLIENT_SECRET=\"\$VIKUNJA_OIDC_CLIENT_SECRET\" \\"
  echo "    --from-literal=OIDC_ISSUER_URL='https://idp-1.tailf15a72.ts.net/' \\"
  echo "    --from-literal=COOKIE_SECRET=\"\$COOKIE_SECRET\""
  echo ""
  
  read -p "Continue deployment without OIDC? (y/N) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo_error "Deployment cancelled. Please create the OIDC secret first."
    exit 1
  fi
  echo_warn "Continuing without OIDC authentication..."
else
  echo_info "OIDC secret 'vikunja-oidc' found"
fi

# Deploy manifests
echo_info "Deploying Vikunja manifests..."
manifest_file="$current_script_dir/vikunja.yaml"
if [[ ! -f "$manifest_file" ]]; then
  echo_error "Manifest file not found: $manifest_file"
  exit 1
fi

# Replace DOMAIN placeholder and apply
sed "s/DOMAIN/$DOMAIN_NAME/g" "$manifest_file" | kubectl apply -f -

echo_info "Waiting for deployment to be ready..."
kubectl rollout status deployment/vikunja -n vikunja --timeout=120s || true

echo ""
echo_info "Vikunja deployment complete!"
echo ""
echo "Access URLs:"
echo "  - Main: https://vikunja.$DOMAIN_NAME"
echo ""
echo "Useful commands:"
echo "  kubectl get pods -n vikunja"
echo "  kubectl logs -n vikunja deployment/vikunja"
echo ""

# Check if OIDC is configured
if kubectl get secret vikunja-oidc -n vikunja >/dev/null 2>&1; then
  echo_info "OIDC authentication is configured"
  echo "Users can log in with their Tailscale identity"
else
  echo_warn "OIDC authentication is NOT configured"
  echo "Vikunja is running with local authentication only"
  echo "To enable OIDC, create the secret as shown above and restart:"
  echo "  kubectl rollout restart deployment/vikunja -n vikunja"
fi

echo ""
echo_info "Next steps:"
echo "1. If OIDC is enabled, test login at https://vikunja.$DOMAIN_NAME"
echo "2. Create your first project and tasks"
echo "3. Use labels for task states: state:ready, state:blocked, earliest-on:YYYY-MM-DD"
echo ""
