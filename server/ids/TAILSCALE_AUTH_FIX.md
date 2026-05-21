# Tailscale tsidp Auth Key Update

When the `tsidp` pod is failing with `Error` status, check the logs for:

```
ERROR failed to start tsnet server error="tsnet.Up: backend: invalid key: API key does not exist"
```

This indicates the Tailscale auth key has expired or been deleted.

## Option 1: Update Secret Using kubectl patch (Recommended)

### Step 1: Generate a new auth key

1. Go to https://login.tailscale.com/admin/settings/keys
2. Create a new **Reusable, Ephemeral** auth key
3. Copy the key (starts with `tskey-auth-`)

### Step 2: Update the Kubernetes secret

```bash
# Replace 'tskey-auth-YOUR-NEW-KEY' with your actual key
kubectl patch secret tsidp-secret -n applications --type='json' -p='[{"op": "replace", "path": "/data/ts-authkey", "value":"'$(echo -n 'tskey-auth-YOUR-NEW-KEY' | base64)'"}]'
```

### Step 3: Restart the deployment

```bash
kubectl rollout restart deployment/tsidp -n applications
```

### Step 4: Verify the fix

```bash
kubectl get pods -n applications -l app=tsidp
kubectl logs -n applications -l app=tsidp --tail=20
```

The pod should now show `Running` status and logs should show successful authentication.

---

## Option 2: Recreate Secret (Alternative)

If patching fails, delete and recreate the secret:

```bash
# Delete old secret
kubectl delete secret tsidp-secret -n applications

# Create new secret with your auth key
kubectl create secret generic tsidp-secret \
  -n applications \
  --from-literal=ts-authkey='tskey-auth-YOUR-NEW-KEY'

# Restart deployment
kubectl rollout restart deployment/tsidp -n applications
```

---

## Troubleshooting

### Check current pod status
```bash
kubectl get pods -n applications | grep tsidp
```

### View pod logs
```bash
kubectl logs -n applications deployment/tsidp --tail=50
```

### Describe pod for events
```bash
kubectl describe pod -n applications -l app=tsidp
```

### Verify secret exists
```bash
kubectl get secret tsidp-secret -n applications
```

Note: The auth key is stored base64-encoded in the secret. Do not commit the actual key to git - use the `ids.yaml` template with `TZ_TOKEN` placeholder.
