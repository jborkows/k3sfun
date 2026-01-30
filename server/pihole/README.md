# Connection issues:
Check
```bash 
kubectl get ingress --all-namespaces
```

If metallb assigned new IPs use
```bash
kubectl get pods -l app=pihole -o name | xargs -I {}  kubectl exec --stdin --tty {} -- vi /etc/pihole/pihole.toml
```
to change addresses
