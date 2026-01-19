# ARC scale-set setup (k3s)

This summary documents the current ARC (scale-set mode) installation and how to target the runner.

## Prereqs
- Helm v3 installed
- Namespace: `actions-runner-system`
- GitHub auth secret: `github-auth` containing `github_token`

## Install controller (scale-set mode)
```bash
helm install arc \
  oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller \
  -n actions-runner-system --create-namespace
```

## Install runner scale set
```bash
helm install my-runners \
  oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set \
  -n actions-runner-system \
  --set githubConfigUrl="https://github.com/jborkows/k3sfun" \
  --set githubConfigSecret=github-auth \
  --set maxRunners=5 \
  --set minRunners=0
```

## Set runner label (scale-set name)
Scale-set mode does **not** support custom runner labels. The only label is the
runner scale set name. We set it to `k3s-home-runner`:

```bash
helm upgrade my-runners \
  oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set \
  -n actions-runner-system \
  --set githubConfigUrl="https://github.com/jborkows/k3sfun" \
  --set githubConfigSecret=github-auth \
  --set runnerScaleSetName=k3s-home-runner \
  --set maxRunners=5 \
  --set minRunners=0
```

Use this in workflows:
```yaml
runs-on: k3s-home-runner
```

## If the PAT token changes
Update the Kubernetes secret and restart the listener so ARC re-authenticates.

```bash
kubectl create secret generic github-auth \
  -n actions-runner-system \
  --from-literal=github_token=NEW_TOKEN \
  --dry-run=client -o yaml | kubectl apply -f -
```

Then restart the listener:
```bash
kubectl delete autoscalinglistener -n actions-runner-system \
  -l actions.github.com/scale-set-name=k3s-home-runner
```

## Enable job containers (kubernetes mode)
Workflows that use a `container:` section require the runner to be in
Kubernetes container mode.

```bash
helm upgrade my-runners \
  oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set \
  -n actions-runner-system \
  --set githubConfigUrl="https://github.com/jborkows/k3sfun" \
  --set githubConfigSecret=github-auth \
  --set runnerScaleSetName=k3s-home-runner \
  --set maxRunners=5 \
  --set minRunners=0 \
  --set containerMode.type=kubernetes \
  --set containerMode.kubernetesModeWorkVolumeClaim.accessModes={ReadWriteOnce} \
  --set containerMode.kubernetesModeWorkVolumeClaim.storageClassName=local-path \
  --set containerMode.kubernetesModeWorkVolumeClaim.resources.requests.storage=1Gi
```

## Restart listener (refresh)
```bash
kubectl delete autoscalinglistener -n actions-runner-system \
  -l actions.github.com/scale-set-name=k3s-home-runner
```

A new listener is created automatically. New runners will register with the
`k3s-home-runner` label when jobs are queued.

## Verify
```bash
kubectl get autoscalinglistener,ephemeralrunnerset,ephemeralrunner -n actions-runner-system
```

