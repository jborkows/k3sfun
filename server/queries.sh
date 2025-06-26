kubectl port-forward svc/prometheus-operated -n monitoring 9090:9090
kubectl logs -n monitoring -l app=prometheus
kubectl exec -n monitoring  power-exporter-8w52f --it -- bash
kubectl get pods -l app=pihole -o jsonpath={..metadata.name}




