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

// assertProductCount checks that ListProducts returns expected number of products.
func assertProductCount(t *testing.T, r *Repo, ctx context.Context, filter products.ProductFilter, expected int) []products.Product {
	t.Helper()
	list, err := r.ListProducts(ctx, filter)
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

// assertProductNotInList checks that a product with given ID is not in the list.
func assertProductNotInList(t *testing.T, list []products.Product, id products.ProductID) {
	t.Helper()
	for _, p := range list {
		if p.ID == id {
			t.Errorf("product ID=%d should not be in list", id)
			return
		}
	}
}

// assertProductMissing checks that a product is in the missing/low list and has Missing=true.
func assertProductMissing(t *testing.T, r *Repo, ctx context.Context, id products.ProductID) {
	t.Helper()
	list, err := r.ListProducts(ctx, products.ProductFilter{OnlyMissingOrLow: true, Limit: 100})
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
func assertProductNotMissing(t *testing.T, r *Repo, ctx context.Context, id products.ProductID) {
	t.Helper()
	list, err := r.ListProducts(ctx, products.ProductFilter{OnlyMissingOrLow: true, Limit: 100})
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

// mustCreateGroup creates a group and fails the test on error.
func mustCreateGroup(t *testing.T, r *Repo, ctx context.Context, name string) products.GroupID {
	t.Helper()
	id, err := r.CreateGroup(ctx, name)
	if err != nil {
		t.Fatalf("CreateGroup(%s): %v", name, err)
	}
	return id
}

// mustCreateProduct creates a product and fails the test on error.
func mustCreateProduct(t *testing.T, r *Repo, ctx context.Context, p products.NewProduct) products.ProductID {
	t.Helper()
	id, err := r.CreateProduct(ctx, p)
	if err != nil {
		t.Fatalf("CreateProduct(%s): %v", p.Name, err)
	}
	return id
}

// mustSetQuantity sets product quantity and fails the test on error.
func mustSetQuantity(t *testing.T, r *Repo, ctx context.Context, id products.ProductID, qty products.Quantity) {
	t.Helper()
	if err := r.SetProductQuantity(ctx, id, qty); err != nil {
		t.Fatalf("SetProductQuantity(%d, %v): %v", id, qty, err)
	}
}

func TestRepo_ListProducts_FilteringAndPaging(t *testing.T) {
	db := openTestDB(t)
	setupCleanDB(t, db)

	r := NewRepo(db)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create test groups with synthetic names
	groupAID := mustCreateGroup(t, r, ctx, "test-group-a")
	groupBID := mustCreateGroup(t, r, ctx, "test-group-b")
	groupCID := mustCreateGroup(t, r, ctx, "test-group-c")

	// Create test products with synthetic names (not real product names)
	testProducts := []products.NewProduct{
		{Name: "test-product-alpha", IconKey: "cart", GroupID: &groupAID, Quantity: 1.5, Unit: products.UnitKG},
		{Name: "test-product-beta", IconKey: "cart", GroupID: &groupAID, Quantity: 2.0, Unit: products.UnitKG},
		{Name: "test-product-gamma", IconKey: "cart", GroupID: &groupBID, Quantity: 0, Unit: products.UnitPiece},
		{Name: "test-product-delta", IconKey: "cart", GroupID: &groupCID, Quantity: 1.0, Unit: products.UnitLiter},
		{Name: "test-śmietanka", IconKey: "cart", GroupID: &groupCID, Quantity: 0, Unit: products.UnitLiter},
	}

	for _, p := range testProducts {
		mustCreateProduct(t, r, ctx, p)
	}

	t.Run("empty filter returns all products", func(t *testing.T) {
		assertProductCount(t, r, ctx, products.ProductFilter{Limit: products.MaxProductsPageSize}, len(testProducts))
	})

	t.Run("filter by single group", func(t *testing.T) {
		list := assertProductCount(t, r, ctx, products.ProductFilter{
			GroupIDs: []products.GroupID{groupAID},
			Limit:    products.MaxProductsPageSize,
		}, 2)
		assertAllInGroup(t, list, groupAID)
	})

	t.Run("filter by multiple groups", func(t *testing.T) {
		assertProductCount(t, r, ctx, products.ProductFilter{
			GroupIDs: []products.GroupID{groupAID, groupBID},
			Limit:    products.MaxProductsPageSize,
		}, 3)
	})

	t.Run("filter by name query", func(t *testing.T) {
		list := assertProductCount(t, r, ctx, products.ProductFilter{
			NameQuery: "alpha",
			Limit:     products.MaxProductsPageSize,
		}, 1)
		assertProductInList(t, list, "test-product-alpha")
	})

	t.Run("filter by name with Polish diacritics", func(t *testing.T) {
		list := assertProductCount(t, r, ctx, products.ProductFilter{
			NameQuery: "Śmie",
			Limit:     products.MaxProductsPageSize,
		}, 1)
		assertProductInList(t, list, "test-śmietanka")
	})

	t.Run("filter by name and group combined", func(t *testing.T) {
		list := assertProductCount(t, r, ctx, products.ProductFilter{
			NameQuery: "delta",
			GroupIDs:  []products.GroupID{groupCID},
			Limit:     products.MaxProductsPageSize,
		}, 1)
		assertProductInList(t, list, "test-product-delta")
	})

	t.Run("filter missing or low quantity", func(t *testing.T) {
		// gamma and śmietanka have quantity=0
		assertProductCount(t, r, ctx, products.ProductFilter{
			OnlyMissingOrLow: true,
			Limit:            products.MaxProductsPageSize,
		}, 2)
	})

	t.Run("count products with filter", func(t *testing.T) {
		count, err := r.CountProducts(ctx, products.ProductFilter{GroupIDs: []products.GroupID{groupAID}})
		if err != nil {
			t.Fatalf("CountProducts: %v", err)
		}
		if count != 2 {
			t.Errorf("expected count=2, got %d", count)
		}
	})

	t.Run("pagination offset", func(t *testing.T) {
		first := assertProductCount(t, r, ctx, products.ProductFilter{Limit: 2, Offset: 0}, 2)
		second := assertProductCount(t, r, ctx, products.ProductFilter{Limit: 2, Offset: 2}, 2)
		assertNoOverlap(t, first, second)
	})

	t.Run("limit is clamped to MaxProductsPageSize", func(t *testing.T) {
		_, err := r.ListProducts(ctx, products.ProductFilter{Limit: 9999, Offset: 0})
		if err != nil {
			t.Fatalf("ListProducts (limit clamp): %v", err)
		}
	})
}

func TestRepo_SetProductQuantity_SyncsMissingFlag(t *testing.T) {
	db := openTestDB(t)
	setupCleanDB(t, db)

	r := NewRepo(db)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create a product with quantity > 0
	id := mustCreateProduct(t, r, ctx, products.NewProduct{
		Name:     "test-qty-sync",
		IconKey:  "cart",
		Quantity: 5,
		Unit:     products.UnitPiece,
	})

	// Set quantity to 0 - should set missing flag
	mustSetQuantity(t, r, ctx, id, 0)
	assertProductMissing(t, r, ctx, id)

	// Set quantity > 0 - should clear missing flag
	mustSetQuantity(t, r, ctx, id, 5)
	assertProductNotMissing(t, r, ctx, id)

	// Set quantity back to 0 - should set missing flag again
	mustSetQuantity(t, r, ctx, id, 0)
	assertProductMissing(t, r, ctx, id)
}

func TestRepo_SuggestProductsByName_PolishDiacritics(t *testing.T) {
	db := openTestDB(t)
	setupCleanDB(t, db)

	r := NewRepo(db)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create products with Polish characters (synthetic test names)
	testProducts := []products.NewProduct{
		{Name: "test-śliwka", IconKey: "cart", Quantity: 1, Unit: products.UnitLiter},
		{Name: "test-żółw", IconKey: "cart", Quantity: 1, Unit: products.UnitKG},
		{Name: "test-ćma", IconKey: "cart", Quantity: 1, Unit: products.UnitPiece},
	}

	for _, p := range testProducts {
		mustCreateProduct(t, r, ctx, p)
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
			suggestions, err := r.SuggestProductsByName(ctx, tc.query, 10)
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
