#!/usr/bin/env bash

#check if pod pi hole with label app=pihole exists 
if kubectl get pod -l "app=pihole" &>/dev/null; then
	echo "Pi-hole namespace already exists."
	exit 0
fi
if [[ -z "$PASS" ]]; then
	echo "Please set the PASS environment variable."
	exit 1
fi

if [[ -z "$DOMAIN_NAME"]]; then
  echo "Please set the DOMAIN_NAME environment variable."
  exit 1
fi
if [[ -z "$CERTIFICATE_NAME"]]; then
  echo "Please set the CERTIFICATE_NAME environment variable."
  exit 1
fi
current_script_dir=$(dirname "$(readlink -f "$0")")
for file in $current_script_dir/*.yaml; do
	echo "Processing file: $file"
sed "s/DOMAIN_NAME/$DOMAIN_NAME/g" *.yaml| \
  sed "s/CERTIFICATE_NAME/$CERTIFICATE_NAME/g" | \
  sed "s/PASS/$PASS/g" > /tmp/$(basename "$file")
  cat "/tmp/$(basename "$file")" | echo
  kubectl apply -f "/tmp/$(basename "$file")"
done
