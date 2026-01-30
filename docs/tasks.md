# Quick Task Reference

This guide maps common tasks to the specific files you need to edit.

---

## Adding SQL Queries

**Primary File**: [internal/db/queries.sql](../internal/db/queries.sql)

1. Add your SQL query to `queries.sql` with `-- name: QueryName :one` or `:many` annotation
2. Run `make sqlc` to regenerate Go code
3. Generated code appears in [internal/db/queries.sql.go](../internal/db/queries.sql.go) (DO NOT EDIT)

**Example pattern**:
```sql
-- name: GetProductByName :one
SELECT * FROM products WHERE name = ?;
```

---

## Modifying UI Templates

**Primary Directory**: [internal/web/views/*.templ](../internal/web/views/)

1. Edit the relevant `.templ` file:
   - [products.templ](../internal/web/views/products.templ) - Product inventory page
   - [shopping.templ](../internal/web/views/shopping.templ) - Shopping list page
   - [admin.templ](../internal/web/views/admin.templ) - Admin panel
   - [pages.templ](../internal/web/views/pages.templ) - Base layouts

2. Run `make templ` to regenerate Go code
3. Generated code appears in `*_templ.go` files (DO NOT EDIT)

---

## Adding Products, Groups, or Icons

**Primary Documentation**: [.opencode/agent/shopping-data.md](../.opencode/agent/shopping-data.md)

**Files involved**:
- [migrations/](../migrations/) - SQL migration files (up/down pairs)
- [web/static/icons/](../web/static/icons/) - SVG icon files

**Process**:
1. Check existing migrations for next number
2. Create SVG icons in `web/static/icons/`
3. Create up migration with groups, icon rules, products
4. Create down migration to reverse changes

---

## Implementing Domain Business Logic

**Primary Directory**: [internal/domain/](../internal/domain/)

**Products domain**:
- [entities.go](../internal/domain/products/entities.go) - Domain entities
- [service.go](../internal/domain/products/service.go) - Write-side service
- [validation.go](../internal/domain/products/validation.go) - Validation rules
- [ports.go](../internal/domain/products/ports.go) - CQRS interfaces

**Shopping list domain**:
- [types.go](../internal/domain/shoppinglist/types.go) - Domain types
- [service.go](../internal/domain/shoppinglist/service.go) - Business logic
- [validation.go](../internal/domain/shoppinglist/validation.go) - Validation rules

---

## Adding HTTP Handlers

**Primary Directory**: [internal/web/](../internal/web/)

1. Create or edit `handlers_*.go` file
2. Add route in `routes_*.go` file
3. Reference [server.go](../internal/web/server.go) for routing patterns

**Existing handlers**:
- [handlers_products.go](../internal/web/handlers_products.go) - Product CRUD
- [handlers_shopping_list.go](../internal/web/handlers_shopping_list.go) - Shopping list
- [handlers_admin.go](../internal/web/handlers_admin.go) - Admin operations
- [handlers_events.go](../internal/web/handlers_events.go) - SSE events

---

## Modifying Database Schema

**Primary Directory**: [migrations/](../migrations/)

1. Check latest migration number
2. Create `NNNN_description.up.sql` for changes
3. Create `NNNN_description.down.sql` for rollback
4. Migrations run automatically on startup

**Schema reference**: [0001_init.up.sql](../migrations/0001_init.up.sql) contains base schema

---

## Project Setup and Configuration

**Primary Documentation**: [README.md](../README.md)

**Key files**:
- [configs/dev.env.example](../configs/dev.env.example) - Environment variables template
- [.env](../.env) - Local configuration (not committed)

---

## Running Tests

| Test Type | Command | Location |
|-----------|---------|----------|
| All tests | `make test` | Entire project |
| Domain tests | `go test ./internal/domain/...` | [internal/domain/](../internal/domain/) |
| Integration tests | `go test ./internal/infrastructure/...` | [internal/infrastructure/persistence/sqlite/](../internal/infrastructure/persistence/sqlite/) |
| Handler tests | `go test ./internal/web/...` | [internal/web/](../internal/web/) |

---

## Kubernetes Deployment

**Primary Documentation**: [README.md](../README.md) (Kubernetes section)

**Key files**:
- [k8s/shoppinglist.yaml](../k8s/shoppinglist.yaml) - Deployment manifest

**Commands**:
```bash
kubectl apply -f k8s/shoppinglist.yaml
```
