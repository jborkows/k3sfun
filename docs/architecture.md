# Architecture Patterns

This document describes the architectural patterns used in the Shopping List application, with links to concrete implementations.

---

## CQRS (Command Query Responsibility Segregation)

CQRS separates read operations (queries) from write operations (commands).

### Pattern
- **Reads**: From views (`v_products`, `v_groups`) for efficient querying
- **Writes**: To base tables with business logic validation

### Implementation

**Read interface (Queries)**:
```go
// internal/domain/products/ports.go
type Queries interface {
    GetProducts(ctx context.Context, filter ProductFilter) ([]ProductView, error)
    GetProductByID(ctx context.Context, id int64) (ProductView, error)
    // ...
}
```

**Write interface (Repository)**:
```go
// internal/domain/products/repository.go
type Repository interface {
    CreateProduct(ctx context.Context, p Product) error
    UpdateProduct(ctx context.Context, p Product) error
    DeleteProduct(ctx context.Context, id int64) error
}
```

**Concrete implementation**: [internal/infrastructure/persistence/sqlite/repo.go](../internal/infrastructure/persistence/sqlite/repo.go)

**Database views**: [migrations/0001_init.up.sql](../migrations/0001_init.up.sql) (v_products, v_groups)

---

## Hexagonal Architecture / Ports & Adapters

Domain defines ports (interfaces), infrastructure provides adapters (implementations).

### Pattern
- **Domain**: Pure business logic, no external dependencies
- **Ports**: Interfaces defined by domain
- **Adapters**: Concrete implementations in infrastructure layer

### Structure

```
internal/
  domain/
    products/
      ports.go        # Read interfaces (Queries)
      repository.go   # Write interfaces (Repository)
      service.go      # Business logic (uses ports)
  infrastructure/
    persistence/sqlite/
      repo.go         # Implements domain ports
    oidc/
      auth.go         # OIDC adapter
```

**Example port definition**: [internal/domain/products/ports.go](../internal/domain/products/ports.go)

**Example adapter**: [internal/infrastructure/persistence/sqlite/repo.go](../internal/infrastructure/persistence/sqlite/repo.go)

---

## SSE (Server-Sent Events) for Real-time Updates

Event-driven UI updates using SSE with pub/sub pattern.

### Pattern
- **Event hub**: Central pub/sub mechanism
- **Topics**: `shopping-list`, `products-list`
- **Client ID**: Unique ID per browser tab for echo suppression

### Implementation

**Event hub**: [internal/web/events.go](../internal/web/events.go)

**Publishing events**:
```go
eventHub.Publish("products-list", Event{
    Type:    "product-updated",
    Payload: product,
})
```

**Client ID handling**: [internal/web/client_id.go](../internal/web/client_id.go)

**Handler**: [internal/web/handlers_events.go](../internal/web/handlers_events.go)

---

## Thin Controllers

HTTP handlers delegate to domain services; business logic lives in domain.

### Pattern
1. Parse and validate input
2. Call domain service
3. Render response

### Example

```go
// internal/web/handlers_products.go
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
    // 1. Parse input
    id := parseID(r)
    req := parseRequest(r)
    
    // 2. Call domain service
    err := h.productService.Update(r.Context(), id, req)
    if err != nil {
        // Handle error
    }
    
    // 3. Render response
    h.renderProducts(w, r)
}
```

**Reference**: [internal/web/handlers_products.go](../internal/web/handlers_products.go)

---

## Domain-Driven Design Concepts

### Entities vs Value Objects
- **Entities**: Have identity (e.g., Product with ID)
- **Value Objects**: Defined by attributes (e.g., Quantity)

### Aggregates
- **Product aggregate**: Product + its rules
- **Shopping list aggregate**: Shopping items + linking logic

**Entity definitions**: [internal/domain/products/entities.go](../internal/domain/products/entities.go)

---

## Dependency Inversion

Domain depends on abstractions (ports), not concrete implementations.

### Example
```go
// Domain service depends on interface
type Service struct {
    repo Repository  // Interface, not *sqlite.Repo
    queries Queries  // Interface, not *sqlite.Queries
}
```

This allows:
- Easy testing with mocks
- Swapping implementations (e.g., test DB vs SQLite)
- No infrastructure imports in domain

---

## Auto-Icon Resolution Pattern

Products get icons based on name pattern matching.

### Pattern
- **Rules table**: `product_icon_rules` with substring patterns
- **Priority ordering**: Higher priority = checked first
- **Default icon**: "cart" when no match

**Database schema**: [migrations/0001_init.up.sql](../migrations/0001_init.up.sql) (product_icon_rules table)
