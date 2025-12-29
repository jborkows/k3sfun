# Repository Guidelines

## Project Structure

```
server/                     # K3s service deployments
  certmanager/              # TLS certificate management
    certmanager.yaml
    init.sh
  filebrowser-externaldrive/ # File browser for external storage
    external_disk.yaml
    filebrowser.yaml
    init.sh
  graphana/                 # Monitoring dashboards
    graphana.yaml
    init.sh
    power.sh
    power.yaml
  ids/                      # Identity service
    ids.yaml
    init.sh
  mainpage/                 # Landing page deployment
    index.html
    init.sh
    mainpage.svg
    mainpage.yaml
    script.js
    style.css
  metal/                    # MetalLB load balancer
    init.sh
    metal.yaml
  pihole/                   # Pi-hole DNS ad blocker
    init.sh
    pi_hole_deployment.yaml
    pi_hole_service_ingress.yaml
    pi_hole_service.yaml
    pi_hole.yaml
  traefik/                  # Traefik ingress controller
    init.sh
    trafik-config.yaml
  backup_code_to_T7.sh      # Backup scripts
  backup_data_to_T7.sh
  envfile.template
  grafana_login.sh
  init.sh
  login_to_docker.sh
  queries.sh

apps/                       # Custom applications
  powerusage/               # Power usage monitoring exporter
    build.sh
    Dockerfile
    power-exporter.sh

mainpage/                   # Landing page source (dev)
  index.html
  mainpage.svg
  mainpage.yml
  script.js
  style.css

.github/
  workflows/
    cleanup-runs.yml        # Cleanup old workflow runs
    deploy-filebrowser.yml  # Filebrowser deployment
    deploy.yml              # Main deployment workflow
    snyk-security.yml       # Security scanning

.worktrees/                 # Git worktrees for subprojects
  k3sfun-shoppinglist/      # Shopping list app (shoppinglist branch)
```

## Branches

| Branch | Purpose |
|--------|---------|
| `master` | K3s infrastructure configurations |
| `shoppinglist` | Go shopping list web application |

Subprojects use git worktrees to allow parallel development:
```bash
git worktree add .worktrees/k3sfun-shoppinglist shoppinglist
```

## Deployment

Each service in `server/` has an `init.sh` script for deployment:
```bash
cd server/<service>
./init.sh
```

Main cluster initialization:
```bash
./initialize.sh
```

## GitHub Actions

- `deploy.yml` - Deploys services to k3s cluster
- `deploy-filebrowser.yml` - Deploys filebrowser service
- `snyk-security.yml` - Security vulnerability scanning
- `cleanup-runs.yml` - Cleans up old workflow runs

## Security

- Never commit secrets; use `envfile.template` as reference
- `.env` and `envfile` are gitignored
- Secrets are managed via k8s secrets or environment variables
