# Vikunja on k3s with Tailscale OIDC

## Project Overview

This project deploys [Vikunja](https://vikunja.io/) - a self-hosted task management application - on a k3s Kubernetes cluster with authentication provided by Tailscale's OIDC Identity Provider (tsidp).

## Architecture

### Components

| Component | Purpose | Documentation |
|-----------|---------|---------------|
| **Vikunja** | Task management API and frontend | https://vikunja.io/docs/ |
| **Tailscale tsidp** | OIDC Identity Provider for authentication | https://github.com/tailscale/tsidp |
| **oauth2-proxy** | OIDC proxy for Vikunja | https://oauth2-proxy.github.io/ |
| **Auto-Transition Sidecar** | Automates task state changes | Custom implementation |
| **SQLite** | Database for tasks and users | Local file storage |
| **k3s** | Kubernetes distribution | https://k3s.io/ |

### Infrastructure

```
User ──▶ Tailscale VPN ──▶ k3s Cluster ──▶ Vikunja Namespace
                                              │
    ┌─────────────────────────────────────────┼──────────────────┐
    │                                         │                  │
    ▼                                         ▼                  ▼
┌─────────┐    ┌──────────┐    ┌──────────┐ ┌──────────┐   ┌──────────┐
│ Traefik │───▶│  Nginx   │───▶│ oauth2-  │ │ Vikunja  │   │  SQLite  │
│ Ingress │    │  Proxy   │    │  proxy   │ │   API    │   │    DB    │
└─────────┘    └──────────┘    └──────────┘ └──────────┘   └──────────┘
                                                     │
                                               ┌──────────┐
                                               │  Auto-   │
                                               │Transition│
                                               │ Sidecar  │
                                               └──────────┘
```

## Task States and Workflow

### Custom Labels for State Tracking

Since Vikunja has fixed task states, we use **labels** to implement custom workflow:

| State Label | Meaning | Transition Trigger |
|-------------|---------|-------------------|
| `state:ready` | Ready to work on | Default state |
| `state:blocked` | Blocked by other tasks | Blockers completed |
| `state:scheduled` | Scheduled for future | Earliest-on date reached |
| `state:on-halt` | Intentionally paused | Manual only |
| `state:in-progress` | Currently being worked on | Manual only |
| `state:completed` | Task finished | Manual or automatic |

### Earliest On Date

Format: `earliest-on:YYYY-MM-DD`

Example: `earliest-on:2025-03-15`

When the date passes, the auto-transition sidecar automatically moves the task to `state:ready`.

### Task Relations (Built-in)

Vikunja's native **"blocked"** and **"blocking"** relations are used for dependencies:
- When a blocking task is completed, the blocked task can be automatically unblocked

## Documentation Links

### Vikunja
- Main docs: https://vikunja.io/docs/
- Task relations: https://vikunja.io/docs/task-relation-kinds/
- Events and listeners: https://vikunja.io/docs/events-and-listeners/
- Configuration: https://vikunja.io/docs/config-options/
- OpenID Connect: https://vikunja.io/docs/openid/

### Tailscale tsidp
- Repository: https://github.com/tailscale/tsidp
- Running tsidp: https://github.com/tailscale/tsidp#running-tsidp
- Configuration: https://github.com/tailscale/tsidp#tsidp-configuration-options

## Deployment

### Manual
```bash
export DOMAIN_NAME=your-domain.com
./vikunja/init.sh
```

### GitHub Actions
- Push to `vikunja` branch triggers deployment
- Requires `DOMAIN_NAME` and `KUBECONFIG_B64` secrets

## Storage

- **Database**: SQLite at `/mnt/external/vikunja-db/vikunja.db`
- **Files**: Task attachments at `/mnt/external/vikunja-files/`
- **Persistent Volume**: Uses same external disk as filebrowser

## Security

- OIDC authentication via Tailscale
- No anonymous access
- HTTPS via Traefik + Let's Encrypt
- OAuth2 proxy handles session management

## Maintenance

### Backup
```bash
# Backup database
kubectl cp vikunja/vikunja-0:/app/vikunja/db/vikunja.db ./vikunja-backup.db
```

### Updates
```bash
# Update Vikunja image
kubectl set image deployment/vikunja vikunja=vikunja/vikunja:latest -n vikunja
```

## Troubleshooting

### Check Pod Status
```bash
kubectl get pods -n vikunja
kubectl describe pod -n vikunja -l app=vikunja
```

### View Logs
```bash
# Vikunja API logs
kubectl logs -n vikunja deployment/vikunja -c vikunja

# OAuth2 proxy logs
kubectl logs -n vikunja deployment/vikunja -c oauth2-proxy

# Auto-transition logs
kubectl logs -n vikunja deployment/vikunja -c auto-transition
```

### Common Issues

**OIDC Login Fails:**
- Verify tsidp is running: `curl https://idp-1.tailf15a72.ts.net/`
- Check K8s secret exists: `kubectl get secret vikunja-oidc -n vikunja`
- Verify redirect URL in tsidp matches: `https://vikunja.DOMAIN/oauth2/callback`

**Database Locked:**
- SQLite doesn't handle concurrent writes well
- Ensure only one Vikunja replica

**Auto-Transition Not Working:**
- Check sidecar logs
- Verify API token is valid
- Ensure tasks have correct label format
