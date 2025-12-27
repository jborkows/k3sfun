# Shopping List Application

A home supplies tracking web application with shopping list management. Built with Go, SQLite, htmx, and templ. UI is in Polish.

## Project Structure & Module Organization

```
cmd/
  shopping/
    main.go           # Application entrypoint (go run ./cmd/shopping)

internal/
  db/                 # sqlc-generated query code
    db.go             # Database connection helpers
    models.go         # Generated model types (DO NOT EDIT)
    queries.sql       # SQL queries (EDIT THIS for new queries)
    queries.sql.go    # Generated Go code from queries.sql (DO NOT EDIT)

  domain/             # Pure domain code (fast unit tests, no DB/network)
    admin/
      ports.go        # Admin ports (e.g., DB optimize)
    products/
      entities.go     # Product/group domain entities
      errors.go       # Domain-specific errors
      ports.go        # CQRS ports (read/write interfaces)
      repository.go   # Repository port definition
      service.go      # Write-side service implementation
      validation.go   # Validation helpers
      validation_test.go
    shoppinglist/
      errors.go       # Shopping list errors
      ports.go        # Shopping list ports
      service.go      # Shopping list service
      types.go        # Shopping list types
      validation.go   # Validation helpers
      validation_test.go

  infrastructure/     # Adapters and integrations
    config/
      config.go       # Env + .env loading
    logging/
      logging.go      # slog logger setup
    oidc/
      auth.go         # OIDC auth against Tailscale tsidp
    persistence/
      sqlite/
        repo.go       # SQLite repository (implements domain ports)
        repo_integration_test.go

  migrator/
    migrator.go       # Database migration runner
    migrator_test.go

  web/                # HTTP handlers and template wiring (thin controllers)
    views/            # templ templates
      admin.templ     # Admin page template
      admin_templ.go  # Generated (DO NOT EDIT)
      helpers.go      # View helper functions
      pages.templ     # Base page layouts
      pages_templ.go  # Generated (DO NOT EDIT)
      products.templ  # Products page template
      products_templ.go # Generated (DO NOT EDIT)
      shopping.templ  # Shopping list template
      shopping_templ.go # Generated (DO NOT EDIT)
      types.go        # View model types
    admin_component.go
    client_id.go
    errors.go
    events.go
    handlers_admin.go
    handlers_admin_test.go
    handlers_events.go
    handlers_icons.go
    handlers_products.go
    handlers_shopping_list.go
    handlers_suggestions.go
    middleware.go
    products_component.go
    routes_admin.go
    routes_products.go
    server.go
    shopping_component.go

migrations/           # SQL migrations (embedded, run on startup)
  0001_init.up.sql
  0001_init.down.sql
  0002_fruits.up.sql
  0002_fruits.down.sql
  embed.go            # Embeds migration files

web/
  static/             # Static assets (htmx-driven UI)
    icons/            # SVG icons
    app.css
    client_id.js
    products.js

k8s/
  shoppinglist.yaml   # Kubernetes deployment manifest

.github/
  workflows/
    deploy.yml        # CI/CD deployment workflow
    snyk-security.yml # Security scanning
```

## Build, Test, and Development

- `make build`: builds `bin/shopping` (runs code generation first).
- `make run`: runs locally (loads `.env` automatically, runs code generation).
- `make dev`: runs with live reload via `air`.
- `make test`: runs `go test ./...` (runs code generation first).
- `make fmt`: runs `gofmt` on Go sources.
- `make gen`: runs both `sqlc` and `templ` code generation.
- `make sqlc`: regenerates `internal/db/*.go` from `queries.sql`.
- `make templ`: regenerates `*_templ.go` files from `.templ` templates.
- `make docker`: builds Docker image.
- `make clean`: removes `./tmp` and `./bin` directories.
- `make migrate-up`: applies migrations using `migrate` CLI (optional; app also auto-migrates).
- `make migrate-down`: rolls back one migration.

## Coding Style & Naming Conventions

- Go formatting: `gofmt` is required (tabs, standard Go style).
- Packages: prefer domain-first names (`products`, `admin`, `shoppinglist`) over generic "service" packages.
- Generated code (DO NOT EDIT manually, regenerate via tooling):
  - `internal/db/db.go`, `internal/db/models.go`, `internal/db/queries.sql.go` (regenerate via `make sqlc`).
  - `internal/web/views/*_templ.go` (regenerate via `make templ`).
- To add new SQL queries: edit `internal/db/queries.sql`, then run `make sqlc`.
- To add/modify templates: edit `*.templ` files in `internal/web/views/`, then run `make templ`.
- Views use `github.com/a-h/templ` for type-safe HTML templating.

## Testing Guidelines

- Use Go's standard `testing` package.
- Keep domain tests fast and isolated (no DB/network):
  - `internal/domain/products/validation_test.go`
  - `internal/domain/shoppinglist/validation_test.go`
- Integration tests (require DB):
  - `internal/infrastructure/persistence/sqlite/repo_integration_test.go`
- Handler tests:
  - `internal/web/handlers_admin_test.go`
- Migration tests:
  - `internal/migrator/migrator_test.go`
- Run all tests with `make test`.

## Commit & Pull Request Guidelines

- Commit messages: use a short imperative summary (e.g., `Add OIDC callback handling`, `Refactor products domain ports`).
- PRs should include: intent/summary, how to test (`make test`, `make run`), and any config changes (new env vars, migrations).

## Domain Concepts

- **Products**: Tracked supplies with name, quantity, unit (kg/litr/sztuk/gramy), minimum quantity, group, icon. `missing` flag for explicit out-of-stock marking. Auto-icon resolution based on name patterns (`product_icon_rules` table).
- **Groups**: Categories like warzywa (vegetables), owoce (fruits), nabiał (dairy), mięso (meat), etc.
- **Shopping List**: Items to buy, optionally linked to products. When marked done, auto-updates product quantity. Auto-links to existing products by name. Done items cleaned up after 6 hours.

## Architecture Patterns

- **CQRS**: Reads from views (`v_products`, `v_groups`), writes to tables. Separate `Queries` (read) and `Repository` (write) interfaces.
- **Hexagonal/Ports & Adapters**: Domain defines ports, infrastructure provides adapters.
- **SSE for real-time updates**: Event hub with pub/sub (`shopping-list`, `products-list` topics), client ID for echo suppression.
- **Thin controllers**: Handlers parse input, call service, render response. Business logic lives in domain services.

## Key Technologies

- **Go** with `modernc.org/sqlite` (pure Go, CGO-free)
- **sqlc** for type-safe SQL code generation
- **templ** for type-safe HTML templates
- **htmx** with SSE extension for hypermedia-driven UI
- **golang-migrate** for database migrations (embedded at build time)
- **OIDC** via `go-oidc` (Tailscale tsidp)
- **OpenTelemetry** for distributed tracing

## Configuration

Key environment variables (see `configs/dev.env.example`):
- `ADDR`: Server address (default `:8080`)
- `DB_DSN`: SQLite database path
- `BASE_URL`: Application base URL
- `AUTH_DISABLED`: Set to `1` to disable OIDC
- `ADMIN_EMAILS`: Comma-separated admin email list
- `LOG_LEVEL`: Logging level (debug/info/warn/error)
- OIDC: `OIDC_ISSUER`, `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`

## Security & Configuration Tips

- Never commit secrets; `.env` is ignored. Rotate OIDC secrets if exposed.
- Admin endpoints require OIDC claim `admin=true` or `ADMIN_EMAILS` (unless `AUTH_DISABLED=1`).
