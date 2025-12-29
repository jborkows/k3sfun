package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"shopping/internal/domain/products"
	"shopping/internal/migrator"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Ping: %v", err)
	}
	return db
}

// setupCleanDB runs migrations and truncates all seeded data to start with a clean slate.
func setupCleanDB(t *testing.T, db *sql.DB) {
	t.Helper()

	if err := migrator.Up(db); err != nil {
		t.Fatalf("migrator.Up: %v", err)
	}

	// Truncate seeded data to start fresh - get table names from sqlite_master
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name NOT LIKE 'schema_migrations'")
	if err != nil {
		t.Fatalf("query tables: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan table name: %v", err)
		}
		tables = append(tables, name)
	}

	for _, table := range tables {
		if _, err := db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %q", table)); err != nil {
			t.Fatalf("truncate %s: %v", table, err)
		}
	}
}

// testEnv holds service and queries interfaces for integration tests.
type testEnv struct {
	svc     *products.Service
	queries products.Queries
}

// setupTestEnv creates a clean test environment with service and queries.
func setupTestEnv(t *testing.T) (testEnv, context.Context) {
	t.Helper()
	db := openTestDB(t)
	setupCleanDB(t, db)

	repo := NewRepo(db)
	svc := products.NewService(repo)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(cancel)

	return testEnv{svc: svc, queries: repo}, ctx
}

// assertProductCount checks that ListProducts returns expected number of products.
func assertProductCount(t *testing.T, queries products.Queries, ctx context.Context, filter products.ProductFilter, expected int) []products.Product {
	t.Helper()
	list, err := queries.ListProducts(ctx, filter)
	if err != nil {
		t.Fatalf("ListProducts: %v", err)
	}
	if len(list) != expected {
		t.Errorf("expected %d products, got %d", expected, len(list))
	}
	return list
}

// assertProductInList checks that a product with given name exists in the list.
func assertProductInList(t *testing.T, list []products.Product, name string) {
	t.Helper()
	for _, p := range list {
		if p.Name == name {
			return
		}
	}
	t.Errorf("expected product %q in list, not found", name)
}

// assertProductMissing checks that a product is in the missing/low list and has Missing=true.
func assertProductMissing(t *testing.T, queries products.Queries, ctx context.Context, id products.ProductID) {
	t.Helper()
	list, err := queries.ListProducts(ctx, products.ProductFilter{OnlyMissingOrLow: true, Limit: 100})
	if err != nil {
		t.Fatalf("ListProducts: %v", err)
	}
	for _, p := range list {
		if p.ID == id {
			if !p.Missing {
				t.Errorf("product ID=%d should have Missing=true", id)
			}
			return
		}
	}
	t.Errorf("product ID=%d should appear in missing/low list", id)
}

// assertProductNotMissing checks that a product is not in the missing/low list.
func assertProductNotMissing(t *testing.T, queries products.Queries, ctx context.Context, id products.ProductID) {
	t.Helper()
	list, err := queries.ListProducts(ctx, products.ProductFilter{OnlyMissingOrLow: true, Limit: 100})
	if err != nil {
		t.Fatalf("ListProducts: %v", err)
	}
	for _, p := range list {
		if p.ID == id {
			t.Errorf("product ID=%d should not be in missing/low list", id)
			return
		}
	}
}

// assertAllInGroup checks that all products in list belong to the expected group.
func assertAllInGroup(t *testing.T, list []products.Product, groupID products.GroupID) {
	t.Helper()
	for _, p := range list {
		if p.GroupID == nil || *p.GroupID != groupID {
			t.Errorf("product %q should be in group %d", p.Name, groupID)
		}
	}
}

// assertNoOverlap checks that two product lists have no common products by ID.
func assertNoOverlap(t *testing.T, list1, list2 []products.Product) {
	t.Helper()
	for _, p1 := range list1 {
		for _, p2 := range list2 {
			if p1.ID == p2.ID {
				t.Errorf("pagination overlap: product %q appears in both pages", p1.Name)
			}
		}
	}
}

// mustCreateGroup creates a group via service and fails the test on error.
func mustCreateGroup(t *testing.T, svc *products.Service, ctx context.Context, name string) products.GroupID {
	t.Helper()
	id, err := svc.CreateGroup(ctx, name)
	if err != nil {
		t.Fatalf("CreateGroup(%s): %v", name, err)
	}
	return id
}

// mustCreateProduct creates a product via service and fails the test on error.
func mustCreateProduct(t *testing.T, svc *products.Service, ctx context.Context, p products.NewProduct) products.ProductID {
	t.Helper()
	id, err := svc.CreateProduct(ctx, p)
	if err != nil {
		t.Fatalf("CreateProduct(%s): %v", p.Name, err)
	}
	return id
}

// mustSetQuantity sets product quantity via service and fails the test on error.
func mustSetQuantity(t *testing.T, svc *products.Service, ctx context.Context, id products.ProductID, qty products.Quantity) {
	t.Helper()
	if err := svc.SetProductQuantity(ctx, id, qty); err != nil {
		t.Fatalf("SetProductQuantity(%d, %v): %v", id, qty, err)
	}
}

func TestProductsService_ListProducts_FilteringAndPaging(t *testing.T) {
	env, ctx := setupTestEnv(t)

	// Create test groups with synthetic names
	groupAID := mustCreateGroup(t, env.svc, ctx, "test-group-a")
	groupBID := mustCreateGroup(t, env.svc, ctx, "test-group-b")
	groupCID := mustCreateGroup(t, env.svc, ctx, "test-group-c")

	// Create test products with synthetic names (not real product names)
	testProducts := []products.NewProduct{
		{Name: "test-product-alpha", IconKey: "cart", GroupID: &groupAID, Quantity: 1.5, Unit: products.UnitKG},
		{Name: "test-product-beta", IconKey: "cart", GroupID: &groupAID, Quantity: 2.0, Unit: products.UnitKG},
		{Name: "test-product-gamma", IconKey: "cart", GroupID: &groupBID, Quantity: 0, Unit: products.UnitPiece},
		{Name: "test-product-delta", IconKey: "cart", GroupID: &groupCID, Quantity: 1.0, Unit: products.UnitLiter},
		{Name: "test-śmietanka", IconKey: "cart", GroupID: &groupCID, Quantity: 0, Unit: products.UnitLiter},
	}

	for _, p := range testProducts {
		mustCreateProduct(t, env.svc, ctx, p)
	}

	t.Run("empty filter returns all products", func(t *testing.T) {
		assertProductCount(t, env.queries, ctx, products.ProductFilter{Limit: products.MaxProductsPageSize}, len(testProducts))
	})

	t.Run("filter by single group", func(t *testing.T) {
		list := assertProductCount(t, env.queries, ctx, products.ProductFilter{
			GroupIDs: []products.GroupID{groupAID},
			Limit:    products.MaxProductsPageSize,
		}, 2)
		assertAllInGroup(t, list, groupAID)
	})

	t.Run("filter by multiple groups", func(t *testing.T) {
		assertProductCount(t, env.queries, ctx, products.ProductFilter{
			GroupIDs: []products.GroupID{groupAID, groupBID},
			Limit:    products.MaxProductsPageSize,
		}, 3)
	})

	t.Run("filter by name query", func(t *testing.T) {
		list := assertProductCount(t, env.queries, ctx, products.ProductFilter{
			NameQuery: "alpha",
			Limit:     products.MaxProductsPageSize,
		}, 1)
		assertProductInList(t, list, "test-product-alpha")
	})

	t.Run("filter by name with Polish diacritics", func(t *testing.T) {
		list := assertProductCount(t, env.queries, ctx, products.ProductFilter{
			NameQuery: "Śmie",
			Limit:     products.MaxProductsPageSize,
		}, 1)
		assertProductInList(t, list, "test-śmietanka")
	})

	t.Run("filter by name and group combined", func(t *testing.T) {
		list := assertProductCount(t, env.queries, ctx, products.ProductFilter{
			NameQuery: "delta",
			GroupIDs:  []products.GroupID{groupCID},
			Limit:     products.MaxProductsPageSize,
		}, 1)
		assertProductInList(t, list, "test-product-delta")
	})

	t.Run("filter missing or low quantity", func(t *testing.T) {
		// gamma and śmietanka have quantity=0
		assertProductCount(t, env.queries, ctx, products.ProductFilter{
			OnlyMissingOrLow: true,
			Limit:            products.MaxProductsPageSize,
		}, 2)
	})

	t.Run("count products with filter", func(t *testing.T) {
		count, err := env.queries.CountProducts(ctx, products.ProductFilter{GroupIDs: []products.GroupID{groupAID}})
		if err != nil {
			t.Fatalf("CountProducts: %v", err)
		}
		if count != 2 {
			t.Errorf("expected count=2, got %d", count)
		}
	})

	t.Run("pagination offset", func(t *testing.T) {
		first := assertProductCount(t, env.queries, ctx, products.ProductFilter{Limit: 2, Offset: 0}, 2)
		second := assertProductCount(t, env.queries, ctx, products.ProductFilter{Limit: 2, Offset: 2}, 2)
		assertNoOverlap(t, first, second)
	})

	t.Run("limit is clamped to MaxProductsPageSize", func(t *testing.T) {
		_, err := env.queries.ListProducts(ctx, products.ProductFilter{Limit: 9999, Offset: 0})
		if err != nil {
			t.Fatalf("ListProducts (limit clamp): %v", err)
		}
	})
}

func TestProductsService_SetQuantity_SyncsMissingFlag(t *testing.T) {
	env, ctx := setupTestEnv(t)

	// Create a product with quantity > 0
	id := mustCreateProduct(t, env.svc, ctx, products.NewProduct{
		Name:     "test-qty-sync",
		IconKey:  "cart",
		Quantity: 5,
		Unit:     products.UnitPiece,
	})

	// Set quantity to 0 via service - should set missing flag
	mustSetQuantity(t, env.svc, ctx, id, 0)
	assertProductMissing(t, env.queries, ctx, id)

	// Set quantity > 0 via service - should clear missing flag
	mustSetQuantity(t, env.svc, ctx, id, 5)
	assertProductNotMissing(t, env.queries, ctx, id)

	// Set quantity back to 0 - should set missing flag again
	mustSetQuantity(t, env.svc, ctx, id, 0)
	assertProductMissing(t, env.queries, ctx, id)
}

func TestProductsService_SetQuantity_RejectsNegativeValue(t *testing.T) {
	env, ctx := setupTestEnv(t)

	id := mustCreateProduct(t, env.svc, ctx, products.NewProduct{
		Name:     "test-validation",
		IconKey:  "cart",
		Quantity: 5,
		Unit:     products.UnitPiece,
	})

	// Negative quantity should fail
	err := env.svc.SetProductQuantity(ctx, id, -1)
	if err == nil {
		t.Error("expected error for negative quantity, got nil")
	}
}

func TestProductsService_CreateProduct_Validation(t *testing.T) {
	env, ctx := setupTestEnv(t)

	t.Run("empty name is rejected", func(t *testing.T) {
		_, err := env.svc.CreateProduct(ctx, products.NewProduct{
			Name: "",
			Unit: products.UnitPiece,
		})
		if err == nil {
			t.Error("expected error for empty name, got nil")
		}
	})

	t.Run("negative quantity is rejected", func(t *testing.T) {
		_, err := env.svc.CreateProduct(ctx, products.NewProduct{
			Name:     "test-negative",
			Quantity: -1,
			Unit:     products.UnitPiece,
		})
		if err == nil {
			t.Error("expected error for negative quantity, got nil")
		}
	})

	t.Run("default icon is assigned when not specified", func(t *testing.T) {
		id := mustCreateProduct(t, env.svc, ctx, products.NewProduct{
			Name: "test-no-icon",
			Unit: products.UnitPiece,
		})
		list, _ := env.queries.ListProducts(ctx, products.ProductFilter{Limit: 100})
		for _, p := range list {
			if p.ID == id && p.IconKey != "cart" {
				t.Errorf("expected default icon 'cart', got %q", p.IconKey)
			}
		}
	})
}

func TestProductsService_SuggestByName_PolishDiacritics(t *testing.T) {
	env, ctx := setupTestEnv(t)

	// Create products with Polish characters (synthetic test names)
	testProducts := []products.NewProduct{
		{Name: "test-śliwka", IconKey: "cart", Quantity: 1, Unit: products.UnitLiter},
		{Name: "test-żółw", IconKey: "cart", Quantity: 1, Unit: products.UnitKG},
		{Name: "test-ćma", IconKey: "cart", Quantity: 1, Unit: products.UnitPiece},
	}

	for _, p := range testProducts {
		mustCreateProduct(t, env.svc, ctx, p)
	}

	tests := []struct {
		query    string
		expected string
	}{
		{"Śliw", "test-śliwka"},
		{"śliw", "test-śliwka"},
		{"ŚLIW", "test-śliwka"},
		{"Żół", "test-żółw"},
		{"żół", "test-żółw"},
		{"Ćma", "test-ćma"},
		{"ćma", "test-ćma"},
	}

	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			suggestions, err := env.queries.SuggestProductsByName(ctx, tc.query, 10)
			if err != nil {
				t.Fatalf("SuggestProductsByName(%q): %v", tc.query, err)
			}
			if len(suggestions) == 0 {
				t.Errorf("expected suggestion for %q, got none", tc.query)
				return
			}
			if suggestions[0].Name != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, suggestions[0].Name)
			}
		})
	}
}
