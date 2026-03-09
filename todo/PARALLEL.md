# Parallel Task Breakdown

## ✅ COMPLETED - Basic Infrastructure

### Files Created
- ✅ `vikunja/vikunja.yaml` - Core K8s manifests
- ✅ `vikunja/init.sh` - Local deployment script
- ✅ `vikunja/auto-transition.sh` - Sidecar automation script
- ✅ `.github/workflows/deploy.yml` - GitHub Actions workflow
- ✅ `AGENTS.md` - Technical documentation
- ✅ `README.md` - Setup instructions
- ✅ `todo/1.md, 2.md, 3.md` - Task tracking

---

## Human Tasks Required (Parallel Possible)

### Task A: Verify tsidp Configuration (You) ⏳
- [ ] Confirm client_id and client_secret from tsidp
- [ ] Verify redirect URL format: `https://vikunja.DOMAIN/oauth2/callback`
- [ ] Test tsidp endpoint: `curl https://idp-1.tailf15a72.ts.net/.well-known/openid-configuration`

### Task B: Verify External Storage (You) ⏳
- [ ] Confirm `/mnt/external` is mounted on k3s nodes
- [ ] Verify write permissions: `touch /mnt/external/test && rm /mnt/external/test`
- [ ] Check available disk space: `df -h /mnt/external`

### Task C: Create K8s Secrets (You - After A)
```bash
# Create namespace
kubectl create namespace vikunja

# Create OIDC secret
kubectl create secret generic vikunja-oidc \
  --namespace=vikunja \
  --from-literal=OIDC_CLIENT_ID='<from-task-a>' \
  --from-literal=OIDC_CLIENT_SECRET='<from-task-a>' \
  --from-literal=OIDC_ISSUER_URL='https://idp-1.tailf15a72.ts.net/' \
  --from-literal=COOKIE_SECRET='$(openssl rand -base64 32 | tr -d '\n')'
```

### Task D: Initial Deployment Test (You - After B)
```bash
export DOMAIN_NAME=your-domain.com
./vikunja/init.sh
```

---

## Remaining Technical Work

### Next Session: OIDC Integration 🔄
**Depends on:** Task A, C complete

Add to `vikunja/vikunja.yaml`:
- oauth2-proxy container
- Update ingress to route through oauth2-proxy
- Health check configuration

### Following Session: Auto-Transition Integration 🔄
**Depends on:** Task D complete (Vikunja running)

Add to `vikunja/vikunja.yaml`:
- auto-transition sidecar container
- ConfigMap for script
- API token secret reference
- Volume mounts

Also need:
- Generate Vikunja API token (via UI or CLI)
- Create `vikunja-api-token` secret
- Test label transitions

---

## Current Status Summary

```
Infrastructure: ✅ COMPLETE
├── Namespace: vikunja
├── Storage: hostPath at /mnt/external/vikunja-* (same as filebrowser)
├── ConfigMap: vikunja-config
├── Deployment: vikunja (single replica)
├── Service: vikunja (ClusterIP)
└── Ingress: vikunja (Traefik + TLS)

Documentation: ✅ COMPLETE
├── AGENTS.md - Architecture & docs
├── README.md - Setup instructions
└── todo/*.md - Task tracking

Scripts: ✅ COMPLETE
├── init.sh - Local deployment
├── auto-transition.sh - Sidecar logic
└── deploy.yml - GitHub Actions

Authentication: ⏳ WAITING
└── Needs: OIDC secret, oauth2-proxy container

Automation: ⏳ WAITING
└── Needs: API token, sidecar container
```

---

## Next Actions

**You can do in parallel:**
1. Complete Task A (verify tsidp)
2. Complete Task B (verify storage)

**After Task A & B complete:**
3. Complete Task C (create secrets)
4. Complete Task D (initial deployment)

**Then I'll implement:**
5. Add oauth2-proxy container (OIDC auth)
6. Add auto-transition sidecar
7. Full integration test
