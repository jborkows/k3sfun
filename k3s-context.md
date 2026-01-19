# Add kube context for remote k3s

Fetch the kubeconfig from the k3s server, replace `127.0.0.1` with the server
address, and add a named context (`my-server-name`).

```bash
ssh user@REMOTE "sudo cat /etc/rancher/k3s/k3s.yaml" > ~/.kube/k3s-remote.yaml
sed -i 's/127.0.0.1/<REMOTE_IP>/g' ~/.kube/k3s-remote.yaml
KUBECONFIG=~/.kube/k3s-remote.yaml kubectl get nodes

KUBECONFIG=~/.kube/config:~/.kube/k3s-remote.yaml kubectl config view --flatten > /tmp/kubeconfig \
  && mv /tmp/kubeconfig ~/.kube/config
kubectl config get-contexts
kubectl config rename-context <old-context-name> my-server-name
kubectl config use-context my-server-name
```
