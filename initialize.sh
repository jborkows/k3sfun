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
  apt-get install -y curl wget jq
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

function enable_treafic_cloudflare(){

MANIFEST="/var/lib/rancher/k3s/server/manifests/traefik.yaml"
BACKUP="${MANIFEST}.bak.$(date +%F_%T)"
NAMESPACE="kube-system"
SECRET_NAME="cloudflare-secret"

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

cat <<EOF > /var/lib/rancher/k3s/server/manifests/traefik-config.yaml
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: traefik
  namespace: kube-system
spec:
  chart: traefik
  repo: https://helm.traefik.io/traefik
  targetNamespace: kube-system
  set:
    additionalArguments[4]: "--log.level=DEBUG"
    additionalArguments[5]: "--entrypoints.web.address=:80"
    additionalArguments[6]: "--entrypoints.websecure.address=:443"
  valuesContent: |-
    api:
      dashboard: true

    entryPoints:
      web:
        address: ":80"
        http:
          redirections:
            entryPoint:
              to: websecure
              scheme: https
      websecure:
        address: ":443"

    providers:
      kubernetesIngress:
        enabled: true
        publishedService:
          enabled: true
      kubernetesCRD:
        enabled: true
EOF
}
function install_metal(){
  echo "🔧 Installing MetalLB..."
  if [[ -z "$IP_RANGE" ]]; then
    echo "Please set the IP_RANGE environment variable."
    exit 1
  fi
  
  if  kubectl get ns metallb-system &>/dev/null; then
    echo "MetalLB namespace already exists."
    return 0
  fi
  kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.10/config/manifests/metallb-native.yaml
  echo "Waiting for MetalLB to be ready..."
  until kubectl get pods -n metallb-system | grep -q 'Running'; do
    echo "⏳ Waiting for MetalLB pods to be Running..."
    sleep 2
  done
  cat <<EOF > /root/yaml/metal.yaml
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: home-pool
  namespace: metallb-system
spec:
  addresses:
  - $IP_RANGE.123-$IP_RANGE.200
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: l2
  namespace: metallb-system
EOF
  kubectl apply -f /root/yaml/metal.yaml
  echo "MetalLB installed successfully!"
}

function install_certmanager(){
  echo "🔧 Installing cert-manager..."
  kubectl apply -f https://github.com/cert-manager/cert-manager/releases/latest/download/cert-manager.yaml
  until kubectl get pods -n cert-manager | grep -q 'Running'; do
    echo "⏳ Waiting for cert-manager pods to be Running..."
    sleep 2
  done
kubectl create secret generic cloudflare-api-token-secret \
  --from-literal=api-token="$CLOUDFLARE_API_TOKEN" \
  -n cert-manager
cat <<EOF > /root/yaml/cert-manager-clusterissuer.yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-dns
spec:
  acme:
    email: your@email.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-dns-key
    solvers:
    - dns01:
        cloudflare:
          email: your@email.com
          apiTokenSecretRef:
            name: cloudflare-api-token-secret
            key: api-token
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: wildcard-borkowskij-com
  namespace: kube-system
spec:
  secretName: wildcard-borkowskij-com
  issuerRef:
    name: letsencrypt-dns
    kind: ClusterIssuer
  commonName: '*.borkowskij.com'
  dnsNames:
    - '*.borkowskij.com'
    - borkowskij.com

EOF
  kubectl apply -f /root/yaml/cert-manager-clusterissuer.yaml
  echo "cert-manager installed successfully!"
}

function install_pi_hole(){
  echo "🔧 Installing Pi-hole..."
  if [[ -z "$PASS" ]]; then
    echo "Please set the PASS environment variable."
    exit 1
  fi
  
  if kubectl get ns pihole &>/dev/null; then
    echo "Pi-hole namespace already exists."
    return 0
  fi

  cat <<EOF > /root/yaml/defaul-wildcard.yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: wildcard-borkowskij-com
  namespace: default
spec:
  secretName: wildcard-borkowskij-com
  dnsNames:
  - '*.borkowskij.com'
  issuerRef:
    name: letsencrypt-dns
    kind: ClusterIssuer
EOF
kubectl apply -f /root/yaml/defaul-wildcard.yaml
  cat <<EOF > /root/yaml/pihole-config.yaml
# pihole-pv.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pihole-pv
spec:
  capacity:
    storage: 2Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /data/pihole
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pihole-dnsmasq-pvc
  namespace: default
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pihole-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pihole
spec:
  replicas: 1
  selector:
    matchLabels:
      app: pihole
  template:
    metadata:
      labels:
        app: pihole
    spec:
      containers:
      - name: pihole
        image: pihole/pihole:latest
        env:
        - name: TZ
          value: "Europe/Berlin"
        - name: WEBPASSWORD
          value: "$PASS"  # Set your admin password here
        ports:
        - containerPort: 80
        - containerPort: 53
        - containerPort: 443
        volumeMounts:
        - name: pihole-config
          mountPath: /etc/pihole
        - name: dnsmasq-config
          mountPath: /etc/dnsmasq.d
      volumes:
      - name: pihole-config
        persistentVolumeClaim:
          claimName: pihole-pvc
      - name: dnsmasq-config
        persistentVolumeClaim:
          claimName: pihole-dnsmasq-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: pihole
spec:
  selector:
    app: pihole
  ports:
    - name: http
      port: 80
      targetPort: 80
    - name: dns-udp
      port: 53
      protocol: UDP
      targetPort: 53
    - name: dns-tcp
      port: 53
      protocol: TCP
      targetPort: 53
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  name: pihole-lb
spec:
  selector:
    app: pihole
  type: LoadBalancer
  ports:
  - name: dns-udp
    port: 53
    protocol: UDP
    targetPort: 53
  - name: dns-tcp
    port: 53
    protocol: TCP
    targetPort: 53
---
apiVersion: v1
kind: Service
metadata:
  name: pihole-web
spec:
  ports:
  - nodePort: 31126
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: pihole
  type: LoadBalancer

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: pihole-ingress
  namespace: default
  annotations:
    traefik.ingress.kubernetes.io/router.entrypoints: websecure
spec:
  rules:
  - host: pihole.borkowskij.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: pihole
            port:
              number: 80
  tls:
  - hosts:
    - pihole.borkowskij.com
    secretName: wildcard-borkowskij-com
EOF
kubectl apply -f /root/yaml/pihole-config.yaml
}

echo "$PASS" | sudo -S -i
echo "$PASS" | sudo -S bash -c "$(declare -f install_utils);   install_utils"
check_k3s
wait_on_k3s
echo "$PASS" | sudo -S bash -c "$(declare -f enable_treafic_cloudflare ); CF_TOKEN='$CF_TOKEN' EMAIL='$EMAIL' enable_treafic_cloudflare "
echo "$PASS" | sudo -S bash -c "$(declare -f install_metal ); IP_RANGE='$IP_RANGE' install_metal" 
echo "$PASS" | sudo -S bash -c "$(declare -f install_certmanager ); CLOUDFLARE_API_TOKEN='$CLOUDFLARE_API_TOKEN' install_certmanager" 
echo "$PASS" | sudo -S bash -c "$(declare -f install_pi_hole ); PASS='$PASS' install_pi_hole" 


