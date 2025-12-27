# Shopping List

A home supplies tracking web application with shopping list management. Track what's missing, current inventory levels, and organize items by groups (e.g., carrots -> vegetables). Built with Go, SQLite, htmx, and templ. UI is in Polish.

## Features

- **Product Inventory**: Track supplies with quantities, units (kg/litr/sztuk/gramy), minimum thresholds, and groups
- **Shopping List**: Add items to buy, mark as done (auto-updates inventory), auto-cleanup after 6 hours
- **Real-time Updates**: SSE-powered UI updates across browser tabs
- **Auto-icon Resolution**: Product names matched against patterns for automatic icon assignment
- **OIDC Authentication**: Secure access via Tailscale tsidp (or disable for development)

## Requirements

- Go 1.25+
- SQLite (via pure-Go `modernc.org/sqlite`, no CGO)
- `sqlc` (for SQL code generation)
- `air` (optional, for live reload)
- `migrate` (golang-migrate CLI, optional - app auto-migrates)

## Quick Start

### 1. Configure

```bash
cp configs/dev.env.example .env
# Edit .env as needed
```

### 2. Run

```bash
# Production-like
make run

# Development with live reload
make dev
```

The app will be available at `http://localhost:8080`.

## Configuration

Key environment variables (see `configs/dev.env.example`):

| Variable | Default | Description |
|----------|---------|-------------|
| `ADDR` | `:8080` | Server listen address |
| `DB_DSN` | `file:data/shopping.db?...` | SQLite database path |
| `BASE_URL` | `http://localhost:8080` | Application base URL |
| `AUTH_DISABLED` | `0` | Set to `1` to skip OIDC auth |
| `ADMIN_EMAILS` | - | Comma-separated admin emails |
| `LOG_LEVEL` | `info` | Logging level (debug/info/warn/error) |
| `OIDC_ISSUER` | - | OIDC provider URL |
| `OIDC_CLIENT_ID` | - | OIDC client ID |
| `OIDC_CLIENT_SECRET` | - | OIDC client secret |
| `OIDC_REDIRECT_URL` | - | OAuth2 callback URL |

## Make Targets

| Command | Description |
|---------|-------------|
| `make build` | Build binary to `bin/shopping` |
| `make run` | Run locally (generates code first) |
| `make dev` | Live reload with `air` |
| `make test` | Run all tests |
| `make fmt` | Format Go code |
| `make gen` | Generate sqlc + templ code |
| `make sqlc` | Generate sqlc code only |
| `make templ` | Generate templ code only |
| `make docker` | Build Docker image |
| `make clean` | Remove `tmp/` and `bin/` |
| `make migrate-up` | Apply migrations manually |
| `make migrate-down` | Rollback one migration |

## Database Migrations

Migrations are embedded and run automatically on startup. Manual usage:

```bash
mkdir -p data
make migrate-up
```

## Architecture

- **CQRS**: Reads from views (`v_products`, `v_groups`), writes to base tables
- **Hexagonal/Ports & Adapters**: Domain defines ports, infrastructure provides adapters
- **SSE Real-time Updates**: Event hub with pub/sub for instant UI updates
- **Thin Controllers**: Handlers delegate to domain services

### Project Structure

```
cmd/shopping/           # Main entrypoint
internal/
  domain/               # Pure business logic (no infrastructure deps)
    products/           # Product/Group entities, service, validation
    shoppinglist/       # Shopping list items, service
    admin/              # Admin operations
  infrastructure/       # Adapters
    persistence/sqlite/ # SQLite repository
    oidc/               # OIDC authentication
    config/             # Environment loading
    logging/            # Structured logging
  db/                   # sqlc-generated code
  migrator/             # Migration runner
  web/                  # HTTP handlers, routes, templates
    views/              # templ templates
migrations/             # SQL migration files
web/static/             # CSS, JS, icons
k8s/                    # Kubernetes manifests
```

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster (e.g., k3s)
- Traefik ingress controller
- Docker Hub credentials (for private images)

### Create OIDC Secret

Create the secret with your OIDC credentials:

```bash
kubectl create secret generic shoppinglist-oidc \
  --namespace=applications \
  --from-literal=OIDC_ISSUER='https://your-open-id-url' \
  --from-literal=OIDC_CLIENT_ID='your-client-id' \
  --from-literal=OIDC_CLIENT_SECRET='your-client-secret' \
  --from-literal=OIDC_REDIRECT_URL='https://shoppinglist.yourdomain.com/oauth2/callback' \
  --from-literal=ADMIN_EMAILS='admin@example.com' \
  --from-literal=DB_DSN='file:data/shopping.db?cache=shared&mode=rwc&_pragma=foreign_keys(1)'
```

To update an existing secret:

```bash
kubectl delete secret shoppinglist-oidc --namespace=applications
# Then recreate with the command above
```

Or use `--dry-run` with `apply`:

```bash
kubectl create secret generic shoppinglist-oidc \
  --namespace=applications \
  --from-literal=OIDC_ISSUER='https://your-open-id-url' \
  --from-literal=OIDC_CLIENT_ID='your-client-id' \
  --from-literal=OIDC_CLIENT_SECRET='your-client-secret' \
  --from-literal=OIDC_REDIRECT_URL='https://shoppinglist.yourdomain.com/oauth2/callback' \
  --from-literal=ADMIN_EMAILS='admin@example.com' \
  --from-literal=DB_DSN='file:data/shopping.db?cache=shared&mode=rwc&_pragma=foreign_keys(1)' \
  --dry-run=client -o yaml | kubectl apply -f -
```

### Create Docker Hub Credentials

```bash
kubectl create secret docker-registry dockerhub-creds \
  --namespace=applications \
  --docker-server=docker.io \
  --docker-username=YOUR_DOCKERHUB_USER \
  --docker-password=YOUR_DOCKERHUB_TOKEN
```

### Deploy

Edit `k8s/shoppinglist.yaml` to set your domain and Docker image, then:

```bash
kubectl apply -f k8s/shoppinglist.yaml
```

### Verify Deployment

```bash
# Check pods
kubectl get pods -n applications -l app=shoppinglist

# Check logs
kubectl logs -n applications -l app=shoppinglist

# Check secret (keys only)
kubectl get secret shoppinglist-oidc -n applications -o jsonpath='{.data}' | jq 'keys'
```

## Admin

- Access `/admin` for admin panel
- `POST /admin/db/optimize` runs `PRAGMA optimize`
- Admin access requires OIDC claim `admin=true` or email in `ADMIN_EMAILS`
- Set `AUTH_DISABLED=1` to bypass authentication (development only)

## Development

### Adding SQL Queries

1. Edit `internal/db/queries.sql`
2. Run `make sqlc`

### Modifying Templates

1. Edit `internal/web/views/*.templ`
2. Run `make templ`

### Testing

```bash
make test
```

- Domain tests: `internal/domain/*/validation_test.go` (fast, no DB)
- Integration tests: `internal/infrastructure/persistence/sqlite/repo_integration_test.go`
- Handler tests: `internal/web/handlers_admin_test.go`
