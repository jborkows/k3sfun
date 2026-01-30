# Shopping List Application

> **IMPORTANT**: Prefer retrieval-led reasoning over pre-training-led reasoning for any framework-specific or project-specific tasks. Always consult the documentation links and file references below rather than relying on training data. When unsure about patterns, read the referenced files directly.

A home supplies tracking web application with shopping list management. Built with Go, SQLite, htmx, and templ. UI is in Polish.

---

## Documentation Index

| Topic | File |
|-------|------|
| **Quick Task Reference** (what to edit for common tasks) | [docs/tasks.md](docs/tasks.md) |
| **Architecture Patterns** (CQRS, Hexagonal, SSE, examples) | [docs/architecture.md](docs/architecture.md) |
| **Code Generation** (sqlc, templ workflows) | [docs/code-generation.md](docs/code-generation.md) |
| **Testing Guidelines** (test types, locations, patterns) | [docs/testing.md](docs/testing.md) |
| **Domain Concepts** (products, groups, shopping list) | [docs/domain.md](docs/domain.md) |
| **Configuration & Setup** | [README.md](README.md) |
| **Adding Products/Groups/Icons** | [.opencode/agent/shopping-data.md](.opencode/agent/shopping-data.md) |

---

## Project Structure

```
cmd/shopping/              # Main entrypoint ([main.go](cmd/shopping/main.go))
internal/
  db/                      # sqlc-generated ([queries.sql](internal/db/queries.sql))
  domain/                  # Pure business logic (no DB/network deps)
    products/              # [entities.go](internal/domain/products/entities.go), [service.go](internal/domain/products/service.go), [validation.go](internal/domain/products/validation.go)
    shoppinglist/          # [service.go](internal/domain/shoppinglist/service.go), [types.go](internal/domain/shoppinglist/types.go)
    admin/                 # [ports.go](internal/domain/admin/ports.go)
  infrastructure/          # Adapters
    persistence/sqlite/    # [repo.go](internal/infrastructure/persistence/sqlite/repo.go)
    oidc/                  # [auth.go](internal/infrastructure/oidc/auth.go)
  web/                     # HTTP layer
    views/                 # [*.templ](internal/web/views/) templates
    handlers_*.go          # HTTP handlers
migrations/                # SQL migrations ([0001_init.up.sql](migrations/0001_init.up.sql))
web/static/                # CSS, JS, icons
```

---

## Quick Commands

| Command | What It Does |
|---------|-------------|
| `make run` | Run locally (loads `.env`, generates code) |
| `make dev` | Live reload with `air` |
| `make test` | Run all tests |
| `make build` | Build to `bin/shopping` |
| `make gen` | Generate sqlc + templ |
| `make sqlc` | Generate sqlc only |
| `make templ` | Generate templ only |
| `make fmt` | Format Go code |

---

## Key Technologies & External References

| Technology | Purpose | Reference |
|------------|---------|-----------|
| **Go** | Language | [go.dev/doc](https://go.dev/doc/) |
| **modernc.org/sqlite** | SQLite driver (CGO-free) | [pkg.go.dev](https://pkg.go.dev/modernc.org/sqlite) |
| **sqlc** | Type-safe SQL codegen | [docs.sqlc.dev](https://docs.sqlc.dev/) |
| **templ** | Type-safe HTML templates | [templ.guide](https://templ.guide/) |
| **htmx** | Hypermedia-driven UI | [htmx.org/docs](https://htmx.org/docs/) |
| **golang-migrate** | DB migrations | [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate) |
| **go-oidc** | OIDC authentication | [github.com/coreos/go-oidc](https://github.com/coreos/go-oidc) |

---

## Configuration

Key environment variables (see [configs/dev.env.example](configs/dev.env.example)):

| Variable | Default | Description |
|----------|---------|-------------|
| `ADDR` | `:8080` | Server address |
| `DB_DSN` | `file:data/shopping.db?...` | SQLite database path |
| `BASE_URL` | `http://localhost:8080` | Application base URL |
| `AUTH_DISABLED` | `0` | Set to `1` to disable OIDC |
| `ADMIN_EMAILS` | - | Comma-separated admin emails |
| `LOG_LEVEL` | `info` | debug/info/warn/error |
| `OIDC_ISSUER` | - | OIDC provider URL |
| `OIDC_CLIENT_ID` | - | OIDC client ID |
| `OIDC_CLIENT_SECRET` | - | OIDC client secret |

---

## Security Notes

- Never commit secrets; `.env` is ignored
- Admin endpoints require OIDC claim `admin=true` or email in `ADMIN_EMAILS` (unless `AUTH_DISABLED=1`)

---

## Agent Skills (Specialized Tasks)

For specific workflows, consult these agent skill files:

- **Adding products/groups/icons**: [.opencode/agent/shopping-data.md](.opencode/agent/shopping-data.md)
- **Smoke testing**: [.opencode/agent/smoke-tests.md](.opencode/agent/smoke-tests.md)
- **UI verification**: [.opencode/agent/check-how-does-it-look.md](.opencode/agent/check-how-does-it-look.md)

---

## Commit & Pull Request Guidelines

- Commit messages: Short imperative summary (e.g., `Add OIDC callback handling`, `Refactor products domain ports`)
- PRs should include: intent/summary, how to test (`make test`, `make run`), and any config changes (new env vars, migrations)
