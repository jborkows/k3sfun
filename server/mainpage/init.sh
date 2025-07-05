kubectl() {
  KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}" command kubectl "$@"
}

if [[ -z "$DOMAIN_NAME" ]]; then
  echo "Please set the DOMAIN_NAME environment variable."
  exit 1
fi
mkdir -p /mnt/external/mainpage

current_script_dir=$(dirname "$(readlink -f "$0")")
(cd $current_script_dir

sed "s/DOMAIN/$DOMAIN_NAME/g" mainpage.yaml > /tmp/mainpage.yaml
sed "s/DOMAIN/$DOMAIN_NAME/g" index.html > /mnt/external/mainpage/index.html
cp *.css /mnt/external/mainpage/
cp *.js /mnt/external/mainpage/
cp *.svg /mnt/external/mainpage/
kubectl apply -f /tmp/mainpage.yaml
)
