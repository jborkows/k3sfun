# Testing Guidelines

This document describes the testing strategy, test locations, and patterns used in the project.

---

## Test Strategy Overview

| Test Type | Location | Characteristics | Speed |
|-----------|----------|-----------------|-------|
| **Domain tests** | [internal/domain/*/](../internal/domain/) | Pure Go, no dependencies | Fast (< 100ms) |
| **Integration tests** | [internal/infrastructure/persistence/sqlite/](../internal/infrastructure/persistence/sqlite/) | Requires SQLite DB | Medium (~1s) |
| **Handler tests** | [internal/web/](../internal/web/) | HTTP testing with mocks | Fast (< 500ms) |
| **Migration tests** | [internal/migrator/](../internal/migrator/) | Database migrations | Medium (~2s) |

---

## Domain Tests

Pure business logic tests with no external dependencies (no DB, no network).

### Location
- [internal/domain/products/validation_test.go](../internal/domain/products/validation_test.go)
- [internal/domain/shoppinglist/validation_test.go](../internal/domain/shoppinglist/validation_test.go)

### Pattern
```go
func TestValidateProductName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid name", "Milk", false},
        {"empty name", "", true},
        {"too long", strings.Repeat("a", 101), true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateProductName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateProductName() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### What to Test
- Validation rules
- Business logic calculations
- Domain entity behavior
- Error conditions

---

## Integration Tests

Tests that verify database interactions using real SQLite.

### Location
- [internal/infrastructure/persistence/sqlite/repo_integration_test.go](../internal/infrastructure/persistence/sqlite/repo_integration_test.go)

### Pattern
```go
func TestRepository_CreateProduct(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewRepository(db)
    
    // Execute
    err := repo.CreateProduct(ctx, product)
    
    // Verify
    if err != nil {
        t.Fatalf("CreateProduct() error = %v", err)
    }
    
    // Query back and verify
    got, err := repo.GetProductByID(ctx, product.ID)
    // ... assertions
}
```

### Setup Pattern
```go
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    
    // Use :memory: or temporary file
    db, err := sql.Open("sqlite", "file::memory:?cache=shared")
    if err != nil {
        t.Fatal(err)
    }
    
    // Run migrations
    if err := runMigrations(db); err != nil {
        t.Fatal(err)
    }
    
    return db
}
```

---

## Handler Tests

HTTP handler tests using `httptest` and mocks.

### Location
- [internal/web/handlers_admin_test.go](../internal/web/handlers_admin_test.go)

### Pattern
```go
func TestHandler_AdminOptimize(t *testing.T) {
    // Setup mocks
    mockAdmin := &mockAdminService{}
    handler := NewHandler(mockAdmin, ...)
    
    // Create request
    req := httptest.NewRequest("POST", "/admin/db/optimize", nil)
    rr := httptest.NewRecorder()
    
    // Execute
    handler.AdminOptimize(rr, req)
    
    // Verify
    if rr.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
    }
}
```

### Mock Pattern
```go
type mockProductService struct {
    products []Product
    err      error
}

func (m *mockProductService) GetProducts(ctx context.Context, filter ProductFilter) ([]Product, error) {
    return m.products, m.err
}
```

---

## Migration Tests

Verify database migrations work correctly.

### Location
- [internal/migrator/migrator_test.go](../internal/migrator/migrator_test.go)

### Pattern
```go
func TestMigrator_Up(t *testing.T) {
    db := setupTestDB(t)
    
    m := NewMigrator(db, embed.FS)
    
    err := m.Up()
    if err != nil {
        t.Fatalf("Up() error = %v", err)
    }
    
    // Verify schema by querying
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
    // ... assertions
}
```

---

## Running Tests

### All Tests
```bash
make test
```

### Specific Packages
```bash
# Domain tests only
go test ./internal/domain/...

# Integration tests
go test ./internal/infrastructure/persistence/sqlite/...

# Handler tests
go test ./internal/web/...

# Migration tests
go test ./internal/migrator/...
```

### With Coverage
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Verbose Output
```bash
go test -v ./internal/domain/products/...
```

---

## Test Naming Conventions

| Pattern | Purpose | Example |
|---------|---------|---------|
| `TestXxx` | Unit test | `TestValidateProductName` |
| `TestXxx_Yyy` | Sub-test | `TestValidateProductName_Empty` |
| `TestIntegration_Xxx` | Integration test | `TestIntegration_RepoCreate` |
| `BenchmarkXxx` | Performance test | `BenchmarkProductList` |

---

## Best Practices

1. **Keep domain tests fast** - No external dependencies
2. **Use table-driven tests** - For multiple test cases
3. **Clean up resources** - `defer db.Close()`, `t.Cleanup()`
4. **Use t.Parallel()** - For independent tests
5. **Mock at port boundaries** - Test domain with mock repository
6. **Integration tests use real DB** - Verify actual SQL works
7. **Test error cases** - Don't just test happy path
8. **Use testify** - Optional: `github.com/stretchr/testify`

---

## Test Data

### Fixtures
Store test data in `testdata/` directories:
```
internal/
  web/
    testdata/
      products.json
      shopping_list.json
```

### Factory Functions
```go
func newTestProduct() Product {
    return Product{
        Name:     "Test Product",
        Quantity: 5,
        Unit:     "sztuk",
    }
}
```
