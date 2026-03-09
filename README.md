# Vikunja Task Manager on k3s

This repository contains Kubernetes manifests and automation for deploying [Vikunja](https://vikunja.io/) task manager on a k3s cluster with Tailscale OIDC authentication.

## Features

- **Vikunja Task Manager** - Self-hosted task management with projects, kanban boards, and CalDAV
- **Tailscale OIDC** - Authentication via Tailscale Identity Provider (tsidp)
- **SQLite Database** - Persistent storage on external disk
- **Auto-Transition** - Automatic task state management based on blockers and dates
- **GitHub Actions** - Automated deployment on push

## Quick Start

### Prerequisites

1. k3s cluster with Traefik ingress
2. Tailscale tsidp running (https://idp-1.tailf15a72.ts.net/)
3. External disk mounted at `/mnt/external`
4. GitHub repository secrets configured

### Setup

1. **Register Vikunja in tsidp:**
   - Go to your tsidp admin UI
   - Register new OIDC client:
     - Client Name: `vikunja`
     - Redirect URI: `https://vikunja.YOUR_DOMAIN/oauth2/callback`
   - Save the generated `client_id` and `client_secret`

2. **Create K8s Secret:**
   ```bash
   kubectl create namespace vikunja
   kubectl create secret generic vikunja-oidc \
     --namespace=vikunja \
     --from-literal=OIDC_CLIENT_ID='<tsidp-client-id>' \
     --from-literal=OIDC_CLIENT_SECRET='<tsidp-secret>' \
     --from-literal=OIDC_ISSUER_URL='https://idp-1.tailf15a72.ts.net/' \
     --from-literal=COOKIE_SECRET='$(openssl rand -base64 32 | tr -d '\n')'
   ```

3. **Deploy:**
   ```bash
   export DOMAIN_NAME=your-domain.com
   ./vikunja/init.sh
   ```
   Or push to trigger GitHub Actions deployment.

## Task States

Vikunja uses labels to track custom states:

| Label | Description |
|-------|-------------|
| `state:ready` | Ready to work on |
| `state:blocked` | Blocked by other task(s) |
| `state:scheduled` | Has earliest-on date |
| `state:on-halt` | Paused/halted |
| `state:in-progress` | Currently being worked on |
| `state:completed` | Done |

### Earliest On

To schedule a task for future activation:
1. Add label: `earliest-on:2025-03-15` (ISO date format)
2. When the date arrives, the auto-transition sidecar moves it to `state:ready`

### Blocked Tasks

Use Vikunja's built-in **"blocked" task relation** to mark dependencies:
1. Edit task → Relations → Add relation
2. Select "blocked by" another task
3. When the blocking task is completed, the auto-transition sidecar moves this task to `state:ready`

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        k3s Cluster                          │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │  Namespace: vikunja                                   │ │
│  │                                                       │ │
│  │  ┌──────────────┐    ┌──────────────┐                │ │
│  │  │   Nginx      │───▶│ oauth2-proxy │──▶ tsidp       │ │
│  │  │   (proxy)    │    │  (OIDC)      │   (Tailscale)  │ │
│  │  └──────────────┘    └──────────────┘                │ │
│  │         │                                            │ │
│  │         ▼                                            │ │
│  │  ┌──────────────┐    ┌──────────────┐                │ │
│  │  │   Vikunja    │◄───│ Auto-Transition              │ │
│  │  │   (API/UI)   │    │   (sidecar)  │                │ │
│  │  └──────────────┘    └──────────────┘                │ │
│  │         │                                            │ │
│  │         ▼                                            │ │
│  │  ┌──────────────┐    ┌──────────────┐                │ │
│  │  │  SQLite DB   │    │   Files      │                │ │
│  │  │ (/mnt/ext)   │    │ (attachments)│                │ │
│  │  └──────────────┘    └──────────────┘                │ │
│  │                                                       │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Documentation

- [Vikunja Documentation](https://vikunja.io/docs/)
- [Task Relations](https://vikunja.io/docs/task-relation-kinds/)
- [Events and Listeners](https://vikunja.io/docs/events-and-listeners/)
- [Tailscale tsidp](https://github.com/tailscale/tsidp)

## GitHub Actions

Deployment is automated via GitHub Actions:
- **Trigger**: Push to `vikunja` branch or tag `vikunja-*`
- **Runner**: `k3s-home-runner`
- **Required Secrets**:
  - `DOMAIN_NAME` - Your domain name
  - `KUBECONFIG_B64` - Base64 encoded kubeconfig

## Local Development

```bash
# Test deployment locally
export DOMAIN_NAME=example.com
./vikunja/init.sh

# View logs
kubectl logs -n vikunja deployment/vikunja -c vikunja
kubectl logs -n vikunja deployment/vikunja -c oauth2-proxy
kubectl logs -n vikunja deployment/vikunja -c auto-transition
```

## License

MIT - See Vikunja license at https://github.com/go-vikunja/vikunja
