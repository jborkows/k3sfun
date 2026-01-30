# Code Generation

This document describes the code generation workflows for sqlc (SQL) and templ (HTML templates).

---

## sqlc - Type-safe SQL

sqlc generates Go code from SQL queries, providing compile-time safety for database operations.

### Workflow

```
internal/db/queries.sql  --(make sqlc)-->  internal/db/queries.sql.go
       (EDIT THIS)                             (DO NOT EDIT)
```

### Adding New Queries

1. **Edit** [internal/db/queries.sql](../internal/db/queries.sql)

2. **Add query with annotation**:
```sql
-- name: GetProductsByGroup :many
SELECT * FROM v_products WHERE group_name = ? ORDER BY name;

-- name: GetProductByID :one
SELECT * FROM v_products WHERE id = ?;

-- name: CreateProduct :exec
INSERT INTO products (name, group_id, quantity_value, quantity_unit, min_quantity_value)
VALUES (?, ?, ?, ?, ?);
```

**Annotation options**:
- `:one` - Returns single row or error
- `:many` - Returns slice of rows
- `:exec` - No return value (INSERT/UPDATE/DELETE)
- `:execresult` - Returns sql.Result

3. **Generate**:
```bash
make sqlc
```

4. **Use generated code**:
```go
products, err := queries.GetProductsByGroup(ctx, "warzywa")
```

### Generated Files (DO NOT EDIT)

| File | Purpose |
|------|---------|
| [internal/db/queries.sql.go](../internal/db/queries.sql.go) | Generated query methods |
| [internal/db/models.go](../internal/db/models.go) | Generated struct types |
| [internal/db/db.go](../internal/db/db.go) | Database connection helpers |

---

## templ - Type-safe HTML Templates

templ generates Go code from HTML templates, providing compile-time template validation.

### Workflow

```
internal/web/views/*.templ  --(make templ)-->  internal/web/views/*_templ.go
         (EDIT THESE)                             (DO NOT EDIT)
```

### Modifying Templates

1. **Edit the `.templ` file**:
   - [products.templ](../internal/web/views/products.templ) - Product inventory UI
   - [shopping.templ](../internal/web/views/shopping.templ) - Shopping list UI
   - [admin.templ](../internal/web/views/admin.templ) - Admin panel
   - [pages.templ](../internal/web/views/pages.templ) - Base page layouts

2. **Template syntax**:
```templ
templ ProductList(products []ProductView) {
    <div class="product-list">
        for _, p := range products {
            @ProductCard(p)
        }
    </div>
}

templ ProductCard(p ProductView) {
    <div class="product-card" data-id={ fmt.Sprintf("%d", p.ID) }>
        <h3>{ p.Name }</h3>
        <p>Quantity: { fmt.Sprintf("%.2f %s", p.Quantity, p.Unit) }</p>
    </div>
}
```

3. **Generate**:
```bash
make templ
```

4. **Use in handlers**:
```go
views.ProductList(products).Render(r.Context(), w)
```

### htmx Integration

Templates use htmx attributes for hypermedia-driven interactions:

```templ
<button hx-post="/products/{p.ID}/add"
        hx-target="#product-list"
        hx-swap="outerHTML">
    Add to Shopping List
</button>
```

### Generated Files (DO NOT EDIT)

| File | Source |
|------|--------|
| [products_templ.go](../internal/web/views/products_templ.go) | products.templ |
| [shopping_templ.go](../internal/web/views/shopping_templ.go) | shopping.templ |
| [admin_templ.go](../internal/web/views/admin_templ.go) | admin.templ |
| [pages_templ.go](../internal/web/views/pages_templ.go) | pages.templ |

---

## Make Commands Reference

| Command | Runs |
|---------|------|
| `make gen` | Both sqlc + templ |
| `make sqlc` | `sqlc generate` (requires sqlc CLI) |
| `make templ` | `templ generate` (requires templ CLI) |

---

## Troubleshooting

### sqlc issues
- **Error**: "table not found" → Check table names in queries match schema
- **Error**: "column not found" → Check column names in migrations

### templ issues
- **Error**: "undefined: views.Xxx" → Run `make templ` to regenerate
- **Template not updating** → Ensure you ran `make templ` after editing `.templ`

---

## Best Practices

1. **Always regenerate after editing** source files
2. **Don't edit generated files** - changes will be overwritten
3. **Commit generated files** - ensures project builds without tooling
4. **Use type-safe parameters** - leverage Go types in templates/queries
